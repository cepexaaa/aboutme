import java.util.concurrent.atomic.AtomicReference
import kotlin.concurrent.atomics.AtomicIntArray
import kotlin.concurrent.atomics.ExperimentalAtomicApi

/**
 * Int-to-Int hash map with open addressing and linear probes.
 * @author Кубеш Сергей 
 */
class IntIntHashMap {
    private var core = AtomicReference(Core(INITIAL_CAPACITY))

    /**
     * Returns value for the corresponding key or zero if this key is not present.
     *
     * @param key a positive key.
     * @return value for the corresponding or zero if this key is not present.
     * @throws IllegalArgumentException if key is not positive.
     */
    operator fun get(key: Int): Int {
        require(key > 0) { "Key must be positive: $key" }
        return toValue(core.get().getInternal(key))
    }

    /**
     * Changes value for the corresponding key and returns old value or zero if key was not present.
     *
     * @param key   a positive key.
     * @param value a positive value.
     * @return old value or zero if this key was not present.
     * @throws IllegalArgumentException if key or value are not positive, or value is equal to
     * [Integer.MAX_VALUE] which is reserved.
     */
    fun put(key: Int, value: Int): Int {
        require(key > 0) { "Key must be positive: $key" }
        require(isValue(value)) { "Invalid value: $value" }
        return toValue(putAndRehashWhileNeeded(key, value))
    }

    /**
     * Removes value for the corresponding key and returns old value or zero if key was not present.
     *
     * @param key a positive key.
     * @return old value or zero if this key was not present.
     * @throws IllegalArgumentException if key is not positive.
     */
    fun remove(key: Int): Int {
        require(key > 0) { "Key must be positive: $key" }
        return toValue(putAndRehashWhileNeeded(key, DEL_VALUE))
    }

    private fun putAndRehashWhileNeeded(key: Int, value: Int): Int {
        while (true) {
            val curCore = core.get()
            val oldValue = curCore.putInternal(key, value)
            if (oldValue != NEEDS_REHASH) return oldValue
            val nc = curCore.rehash()
            core.compareAndSet(curCore, nc)
        }
    }

    @OptIn(ExperimentalAtomicApi::class)
    private class Core(capacity: Int) {
        // Pairs of <key, value> here, the actual
        // size of the map is twice as big.
        val map: AtomicIntArray //: IntArray = IntArray(2 * capacity) //AtomicIntArray(2 * capacity)//
        val shift: Int
        val next = AtomicReference<Core?>(null)

        init {
            val mask = capacity - 1
            assert(mask > 0 && mask and capacity == 0) { "Capacity must be power of 2: $capacity" }
            shift = 32 - Integer.bitCount(mask)
            map = AtomicIntArray(2 * capacity)
        }

        fun getInternal(key: Int): Int {
            var index = index(key)
            var probes = 0
            while (map.loadAt(index) != key) { // optimize for successful lookup
                if (map.loadAt(index) == NULL_KEY) return NULL_VALUE // not found -- no value
                if (++probes >= MAX_PROBES) return NULL_VALUE
                if (index == 0) index = map.size
                index -= 2
            }
            val v = map.loadAt(index + 1)
            if (isMovedValue(v)) {
                val newCore = next.get()!!
                if (v != MOVED_NULL) {
                    newCore.copy2NewTable(key, unmovedValue(v))
                }
                return newCore.getInternal(key)
            }
            return v
//            if (v == DEL_VALUE) return NULL_VALUE
//            return v
        }

        fun putInternal(key: Int, value: Int): Int {
            var index = index(key)
            var probes = 0
            while (true) {
                val curKey = map.loadAt(index)
                if (curKey == key) {
                    return processVal(key, value, index)
                }
                if (curKey == NULL_KEY) {
                    if (value == DEL_VALUE) return NULL_VALUE // если мы хотим удалить значение по у=ключу, которого ещё нет, то ничего не делаем
                    if (map.compareAndSetAt(index, NULL_KEY, key)) {
                        return processVal(key, value, index)
                    }
                    continue
                }
                if (++probes >= MAX_PROBES) return NEEDS_REHASH
                if (index == 0) index = map.size
                index -= 2
            }
//            var index = index(key)
//            var probes = 0
//            while (true) {
//                if (map.compareAndSetAt(index, NULL_KEY, key)) {
//                    if (map.compareAndSetAt(index + 1, NULL_VALUE, value)) {
//                        updateWithRehash(key, value)
//                        return NULL_VALUE
//                    }
//                }
//                val curKey = map.loadAt(index)
//                if (curKey == key) {
//                    val oldValue = map.loadAt(index + 1)
//                    if (!map.compareAndSetAt(index + 1, oldValue, value)) {
//                        continue
//                    }
//                    updateWithRehash(key, value)
//                    return oldValue
//                }
//                if (++probes >= MAX_PROBES) {
//                    return NEEDS_REHASH
//                }
//                if (index == 0) index = map.size
//                index -= 2
//            }
        }

        fun processVal(key: Int, value: Int, index: Int): Int {
            while (true) {
                val oldValue = map.loadAt(index + 1)
                if (isMovedValue(oldValue)) {
                    val newCore = next.get()!!
                    if (oldValue != MOVED_NULL) {
                        newCore.copy2NewTable(key, unmovedValue(oldValue))
                    }
                    return newCore.putInternal(key, value)
                }
                if (map.compareAndSetAt(index + 1, oldValue, value)) {
                    return oldValue
                }
            }
        }

        // раньше добавлял ключ - если получилось, то значение
        // но это не работает
        // просто добавляю всё через cas
        // Проблема - удалённое значение из новой таблицы можно заменить старым.
        // решение удаление делаем только из старой
        // а вставку в новую
        fun rehash(): Core {
            val existingNext = next.get()
            var newCore: Core
            if (existingNext == null) {
                newCore = Core(map.size)
                next.compareAndSet(null, newCore)
            }
            newCore = next.get()!!
//            if (next.compareAndSet(null, newCore)) {

            var index = 0
//            while (index < map.size) {
//                val value = map.loadAt(index + 1)
//                val key = map.loadAt(index)
//                if (key != NULL_KEY || isValue(value)) {
//                    newCore.putRehash(key, value)
//                }
//                index += 2
//            }
//            }
            while (index < map.size) {
                while (true) {
                    val curValue = map.loadAt(index + 1)
                    if (isMovedValue(curValue)) {
                        if (curValue != MOVED_NULL) {
                            val curKey = map.loadAt(index)
                            newCore.copy2NewTable(curKey, unmovedValue(curValue))
                        }
                        break
                    }

                    val movedMarker = if (curValue == NULL_VALUE || curValue == DEL_VALUE) MOVED_NULL else movedValue(curValue)
                    if (map.compareAndSetAt(index + 1, curValue, movedMarker)) {
                        if (isValue(curValue)) {
                            val curKey = map.loadAt(index)
                            newCore.copy2NewTable(curKey, curValue)
                        }
                        break
                    }
                }
                index += 2
            }
            return newCore
        }

        fun putRehash(key: Int, value: Int) {
            val index = index(key)
            map.compareAndSetAt(index, NULL_KEY, key)
            map.compareAndSetAt(index + 1, NULL_VALUE, value)
        }

        fun updateWithRehash(key: Int, value: Int): Int {
            val nextCore = next.get()
            // если многократное рехэширование произошло, то надо while (ture)
            if (nextCore != null) {
                return nextCore.putInternal(key, value)

            }
            return DEL_VALUE
        }

        fun copy2NewTable(key: Int, value: Int) {
            var index = index(key)
            var probes = 0
            while (true) {
                val curKey = map.loadAt(index)
                if (curKey == key) {
                    map.compareAndSetAt(index + 1, NULL_VALUE, value)
                    return
                }
                if (curKey == NULL_KEY) {
                    if (map.compareAndSetAt(index, NULL_KEY, key)) {
                        map.compareAndSetAt(index + 1, NULL_VALUE, value)
                        return
                    }
                    continue
                }
                if (++probes >= MAX_PROBES) {
                    return
                }
                if (index == 0) index = map.size
                index -= 2
            }
        }

        /**
         * Returns an initial index in map to look for a given key.
         */
        fun index(key: Int): Int = (key * MAGIC ushr shift) * 2
    }
}

private fun isMovedValue(value: Int): Boolean = value < 0 && value != NEEDS_REHASH
private fun movedValue(value: Int): Int = -(value + 1)
private fun unmovedValue(value: Int) = movedValue(value)
private const val MOVED_NULL = Int.MIN_VALUE

private const val MAGIC = -0x61c88647 // golden ratio
private const val INITIAL_CAPACITY = 2 // !!! DO NOT CHANGE INITIAL CAPACITY !!!
private const val MAX_PROBES = 8 // max number of probes to find an item
private const val NULL_KEY = 0 // missing key (initial value)
private const val NULL_VALUE = 0 // missing value (initial value)
private const val DEL_VALUE = Int.MAX_VALUE // mark for removed value
private const val NEEDS_REHASH = -1 // returned by `putInternal` to indicate that rehash is needed

// Checks is the value is in the range of allowed values
private fun isValue(value: Int): Boolean = value in (1 until DEL_VALUE)

// Converts internal value to the public results of the methods
private fun toValue(value: Int): Int = if (isValue(value)) value else 0

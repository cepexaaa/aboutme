import kotlin.concurrent.atomics.AtomicReference
import kotlin.concurrent.atomics.ExperimentalAtomicApi

/**
 * AtomicArray with practical lock-free implementation of CAS and CAS2 operations
 * based on paper by Timothy L. Harris, Keir Fraser and Ian A. Pratt:
 * [A Practical Multi-Word Compare-and-Swap Operation](https://www.cl.cam.ac.uk/research/srg/netos/papers/2002-casn.pdf).
 *
 * @author Кубеш Сергей
 */
@Suppress("UNCHECKED_CAST")
@OptIn(ExperimentalAtomicApi::class)
class AtomicArray<E : Any>(size: Int, initialValue: E) {
    private val a: Array<ObjR<E>> = Array(size) { ObjR(initialValue) }

    class ObjR<E : Any>(e: E) {
        val elem: AtomicReference<Any> = AtomicReference(e)

        fun get(): E {
            while (true) {
                when (val cur = elem.load()) {
                    is Task -> {
                        cur.apply()
                    }
                    else -> return cur as E
                }
            }
        }

        fun replace(expect: Any, update: Any): Boolean {
            while (true) {
                when (val cur = elem.load()) {
                    expect -> {
                        if (elem.compareAndSet(expect, update)) return true
                    }
                    is Task -> {
                        cur.apply()
                        continue
                    }
                    else -> return false
                }
            }
        }
    }

    abstract class Task {
        abstract fun apply(): Boolean
    }

    class DCSS<E : Any>(
        val a: ObjR<E>,
        val expectA: Any,
        val updateA: Descriptor<E>
    ) : Task() {

        val status: AtomicReference<Status> = AtomicReference(Status.UNDECIDED)

        override fun apply(): Boolean {
            if (updateA.status.load() == Status.UNDECIDED) {
                status.compareAndSet(Status.UNDECIDED, Status.SUCCESS)
            } else {
                status.compareAndSet(Status.UNDECIDED, Status.FAILED)
            }

            return when (status.load()) {
                Status.SUCCESS -> {
                    a.elem.compareAndSet(this, updateA)
                    true
                }
                Status.FAILED -> {
                    a.elem.compareAndSet(this, expectA)
                    false
                }
                else -> false
            }
        }
    }

    class Descriptor<E : Any>(
        val a1: ObjR<E>, val expect1: E, val update1: E,
        val a2: ObjR<E>, val expect2: E, val update2: E
    ) : Task() {

        val status: AtomicReference<Status> = AtomicReference(Status.UNDECIDED)

        override fun apply(): Boolean {
            if (a2.elem.load() != this) {
                val desc2 = DCSS(a2, expect2 as Any, this)
                if (a2.replace(expect2 as Any, desc2)) {
                    desc2.apply()
                }
            }

            val secondStatus = when (val secondVal = a2.elem.load()) {
                this -> Status.SUCCESS
                is DCSS<*> -> {
                    (secondVal as DCSS<E>).apply()
                    if (a2.elem.load() == this) Status.SUCCESS else Status.FAILED
                }
                else -> Status.FAILED
            }

            if (secondStatus == Status.SUCCESS) {
                status.compareAndSet(Status.UNDECIDED, Status.SUCCESS)
            } else {
                status.compareAndSet(Status.UNDECIDED, Status.FAILED)
            }

            return when (status.load()) {
                Status.SUCCESS -> {
                    a1.elem.compareAndSet(this, update1 as Any)
                    a2.elem.compareAndSet(this, update2 as Any)
                    true
                }
                Status.FAILED -> {
                    a1.elem.compareAndSet(this, expect1 as Any)
                    false
                }
                else -> false
            }
        }
    }

    fun get(index: Int): E = a[index].get()

    fun cas(index: Int, expected: E, update: E): Boolean {
        while (true) {
            when (val cur = a[index].elem.load()) {
                expected -> {
                    if (a[index].elem.compareAndSet(expected, update)) return true
                }
                is Task -> {
                    cur.apply()
                    continue
                }
                else -> return false
            }
        }
    }

    fun cas2(
        index1: Int, expected1: E, update1: E,
        index2: Int, expected2: E, update2: E
    ): Boolean {
        if (index1 == index2) {
            if (expected1 != expected2) return false
            return cas(index1, expected1, update2)
        }

        val descriptor = if (index1 < index2) {
            Descriptor(
                a[index1], expected1, update1,
                a[index2], expected2, update2
            )
        } else {
            Descriptor(
                a[index2], expected2, update2,
                a[index1], expected1, update1
            )
        }

        return if (descriptor.a1.replace(descriptor.expect1 as Any, descriptor)) {
            descriptor.apply()
        } else {
            false
        }
    }

    enum class Status { UNDECIDED, FAILED, SUCCESS }
}
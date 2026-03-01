import java.util.concurrent.*
import java.util.concurrent.atomic.*

/**
 * @author Кубеш Сергей Дмитриевич
 */
class FlatCombiningQueue<E> : Queue<E> {
    private val queue = ArrayDeque<E>() // sequential queue
    private val combinerLock = AtomicBoolean(false) // unlocked initially
    private val tasksForCombiner = AtomicReferenceArray<Any?>(TASKS_FOR_COMBINER_SIZE)

    override fun enqueue(element: E) {
        // TODO: Make this code thread-safe using the flat-combining technique.
        // TODO: 1.  Try to become a combiner by
        // TODO:     changing `combinerLock` from `false` (unlocked) to `true` (locked).
        // TODO: 2a. On success, apply this operation and help others by traversing
        // TODO:     `tasksForCombiner`, performing the announced operations, and
        // TODO:      updating the corresponding cells to `Result`.
        // TODO: 2b. If the lock is already acquired, announce this operation in
        // TODO:     `tasksForCombiner` by replacing a random cell state from
        // TODO:      `null` with the element. Wait until either the cell state
        // TODO:      updates to `Result` (do not forget to clean it in this case),
        // TODO:      or `combinerLock` becomes available to acquire.
        while (true) {
            if (combinerLock.compareAndSet(false, true)) {
                queue.addLast(element)
                helpOther()
                combinerLock.set(false)
                return
            } else {
                val index = findSlotForOperation(EnqueueOp(element))
                if (index != -1) {
                    when (val result = waitForCompletion(index)) {
                        PROCESSED -> return
                        else -> continue
                    }
                }
//                var index = -1
//                for (i in 0 until tasksForCombiner.length()) {
//                    val ind = (randomCellIndex() + i) % tasksForCombiner.length()
//                    if (tasksForCombiner.compareAndSet(ind, null, element)) {
//                        index = ind
//                        break
//                    }
//                }
//                val index = findSlotForOperation(EnqueueOp(element))
//                if (index != -1) {
//                    waitForCompletion(index)?.let {return}
//                }
            }
        }
//        queue.addLast(element)
    }

    override fun dequeue(): E? {
        // TODO: Make this code thread-safe using the flat-combining technique.
        // TODO: 1.  Try to become a combiner by
        // TODO:     changing `combinerLock` from `false` (unlocked) to `true` (locked).
        // TODO: 2a. On success, apply this operation and help others by traversing
        // TODO:     `tasksForCombiner`, performing the announced operations, and
        // TODO:      updating the corresponding cells to `Result`.
        // TODO: 2b. If the lock is already acquired, announce this operation in
        // TODO:     `tasksForCombiner` by replacing a random cell state from
        // TODO:      `null` with `Dequeue`. Wait until either the cell state
        // TODO:      updates to `Result` (do not forget to clean it in this case),
        // TODO:      or `combinerLock` becomes available to acquire.
        while (true) {
            if (combinerLock.compareAndSet(false, true)) {
                val result = queue.removeFirstOrNull()
                helpOther()
                combinerLock.set(false)
                return result
            } else {
                val index = findSlotForOperation(Dequeue)
                if (index != -1) {
                    when (val result = waitForCompletion(index)) {
                        is Result<*> -> {
                            @Suppress("UNCHECKED_CAST")
                            return result.value as E?
                        }
                        else -> continue
                    }
//                    val result = waitForCompletion(index)
//                    if (result is Result<*>) {
//                        @Suppress("UNCHECKED_CAST")
//                        return result.value as E?
//                    } else if (result == PROCESSED) {
//                        continue
//                    }
                }
            }
        }
//        return queue.removeFirstOrNull()
    }

    private fun findSlotForOperation(operation: Any?): Int {
        for (i in 0 until tasksForCombiner.length()) {
            val index = (randomCellIndex() + i) % tasksForCombiner.length()
            if (tasksForCombiner.compareAndSet(index, null, operation)) {
                return index
            }
        }
        return  -1
    }

    private fun waitForCompletion(index: Int): Any? {
        while (true) {
            if (combinerLock.compareAndSet(false, true)) {
                helpOther()
                val current = tasksForCombiner.get(index)
                if (current is Result<*>) {
                    tasksForCombiner.set(index, null)
                    combinerLock.set(false)
                    return current
                }
                val res = performOperation(index)
                combinerLock.set(false)
                return res ?: PROCESSED
//                return res
            }

            val current = tasksForCombiner.get(index)
            if (current is Result<*>) {
                tasksForCombiner.set(index, null)
                return current
            }
            if (current == null || current == PROCESSED) {
                return PROCESSED
            }
        }
    }

    private fun helpOther() {
        for (i in 0 until tasksForCombiner.length()) {
            val operation = tasksForCombiner.get(i)
            if (operation != null && operation != PROCESSED && operation !is Result<*>) {
                performOperation(i)
            }
        }
    }

    private fun performOperation(index: Int): Any? {
        val operation = tasksForCombiner.get(index) ?: return null

        return when (operation) {
            is Dequeue -> {
                val result = queue.removeFirstOrNull()
                tasksForCombiner.set(index, Result(result))
                result
            }
            is EnqueueOp<*> -> {
                @Suppress("UNCHECKED_CAST")
                val element = operation.element as E
                queue.addLast(element)
                tasksForCombiner.set(index, PROCESSED)
                PROCESSED
            }
            else -> {
                tasksForCombiner.set(index, null)
                null
            }
        }

    }

    companion object {
        private val PROCESSED = Any()
    }

    private fun randomCellIndex(): Int =
        ThreadLocalRandom.current().nextInt(tasksForCombiner.length())
}

private const val TASKS_FOR_COMBINER_SIZE = 3 // Do not change this constant!

// TODO: Put this token in `tasksForCombiner` for dequeue().
// TODO: enqueue()-s should put the inserting element.
private object Dequeue
private class EnqueueOp<E>(val element: E)

// TODO: Put the result wrapped with `Result` when the operation in `tasksForCombiner` is processed.
private class Result<V>(
    val value: V
)

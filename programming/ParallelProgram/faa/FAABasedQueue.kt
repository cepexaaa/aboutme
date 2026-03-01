import java.util.concurrent.atomic.*

/**
 * @author Кубеш Сергей
 *
 * TODO: Copy the code from `FAABasedQueueSimplified`
 * TODO: and implement the infinite array on a linked list
 * TODO: of fixed-size `Segment`s.
 */
class FAABasedQueue<E> : Queue<E> {
    private val enqIdx = AtomicLong(0)
    private val deqIdx = AtomicLong(0)
    private var head = AtomicReference<Segment>(Segment(0))
    private val tail = AtomicReference<Segment>(head.get())


    // Плохая реализация. Надо искать не с head а с curTail.
    // Но когда я передаю в findSegment curTail,
    // то передаётся или возвращается как будто просто tail.
    // И enqueue вставляет не туда элемент и dequeue не может его достать.
    // Поэтому приходится передавать начало, чтобы указатель не выбежал куда-то дальше
    override fun enqueue(element: E) {
        while (true) {
//            val curTail = tail
            val i = enqIdx.getAndIncrement()

            val s = findSegment(head, i / SEGMENT_SIZE)
            if (s.cells.compareAndSet(i.toInt() % SEGMENT_SIZE, null, element)) {
                return
            }
        }
    }

    @Suppress("UNCHECKED_CAST")
    override fun dequeue(): E? {
        while (true) {
            if (deqIdx.get() >= enqIdx.get()) return null
            val curHead = head
            val i = deqIdx.getAndIncrement().toInt()
            val s = findSegment(curHead, i / SEGMENT_SIZE.toLong())
//            head = AtomicReference<Segment>(s)
            if (s.cells.compareAndSet(i % SEGMENT_SIZE, null, POISONED)) {
                continue
            }
            val element = s.cells.get(i % SEGMENT_SIZE)
            if (s.cells.compareAndSet(i % SEGMENT_SIZE, element, POISONED)) {
                return element as E
            }
        }
    }

//    private fun enqHelper(): Segment {
//
//    }

    private fun findSegment(start: AtomicReference<Segment>, id: Long): Segment {
        var current = start.get()

        while (current.id < id) {
            var next = current.next.get()
            if (next == null) {
                next = Segment(current.id + 1)
                if (current.next.compareAndSet(null, next)) {
                    tail.compareAndSet(current, next)
                } else {
                    next = current.next.get()!!
                }
            }
            current = next
        }
        return current
    }

}

private class Segment(val id: Long) {
    var next = AtomicReference<Segment?>(null)
    val cells = AtomicReferenceArray<Any?>(SEGMENT_SIZE)
}

// DO NOT CHANGE THIS CONSTANT
private const val SEGMENT_SIZE = 2
private val POISONED = Any()

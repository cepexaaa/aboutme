import java.util.concurrent.atomic.*

/**
 * Implementation of the Michael-Scott queue algorithm.
 *
 * @author Кубеш Сергей
 */
class MichaelScottQueue<E> {
    private val head: AtomicReference<Node<E>>
    private val tail: AtomicReference<Node<E>>

    init {
        val dummy = Node<E>(null)
        head = AtomicReference(dummy)
        tail = AtomicReference(dummy)
    }

    fun enqueue(element: E) {
        val node = Node(element)
        while (true) {
            val curTail = tail.get()
            if (curTail.next.compareAndSet(null, node)) {
                tail.compareAndSet(curTail, node)
                return
            } else {
                tail.compareAndSet(curTail, curTail.next.get())
            }
        }
    }

//    fun dequeue(): E? {
//        while (true) {
//            val curHead = head.get()
//            val curHeadNext = head.get().next.get() ?: return null
//            if (head.compareAndSet(curHead, curHeadNext)) {
//                curHead.element = null
//                return curHeadNext.element
//            }
//        }
//    }
fun dequeue(): E? {
    while (true) {
        val curHead = head.get()
        val firstRealElement = curHead.next.get() ?: return null

        if (head.compareAndSet(curHead, firstRealElement)) {
            val result = firstRealElement.element
            firstRealElement.element = null
            return result
        }
    }
}

    // FOR TEST PURPOSE, DO NOT CHANGE IT.
    fun validate() {
        check(tail.get().next.get() == null) {
            "At the end of the execution, `tail.next` must be `null`"
        }
        check(head.get().element == null) {
            "At the end of the execution, the dummy node shouldn't store an element"
        }
    }

    private class Node<E>(
        var element: E?
    ) {
        val next = AtomicReference<Node<E>?>(null)
    }
}

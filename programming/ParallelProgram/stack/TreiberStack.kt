import java.util.concurrent.atomic.*

/**
 * @author Кубеш Сергей
 */
class TreiberStack<E> : Stack<E> {
    // Initially, the stack is empty.
    private val top = AtomicReference<Node<E>?>(null)

    override fun push(element: E) {
        while (true) {
            val curTop = top.get()
            val node = Node(element, curTop)
            if (top.compareAndSet(curTop, node)) {
                return
            }
        }
    }

    override fun pop(): E? {
        while (true) {
            val curTop = top.get() ?: return null
            val newTop = curTop.next
            if (top.compareAndSet(curTop, newTop)) {
                return curTop.element
            }
        }
    }

    private class Node<E>(
        val element: E,
        val next: Node<E>?
    )
}

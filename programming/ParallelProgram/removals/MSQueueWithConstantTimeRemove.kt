@file:Suppress("DuplicatedCode", "FoldInitializerAndIfToElvis")

import java.util.concurrent.atomic.*

/**
 * Implementation of the Michael-Scott queue algorithm with constant-time removal.
 *
 * @author Кубеш Сергей
 *
 * TODO: Use your MichaelScottQueue implementation as a starting point.
 * TODO: Read instructions in the code.
 */
class MSQueueWithConstantTimeRemove<E> : QueueWithRemove<E> {
    private val head: AtomicReference<Node<E>>
    private val tail: AtomicReference<Node<E>>

    init {
        val dummy = Node<E>(element = null, prev = null)
        head = AtomicReference(dummy)
        tail = AtomicReference(dummy)
    }

    override fun enqueue(element: E) {
        // TODO: When adding a new node, check whether
        // TODO: the previous tail is logically removed.
        // TODO: If so, remove it physically from the linked list.
        val newNode = Node(element, null, tail, head)

        while (true) {
            val curTail = tail.get()
//            val tailNext = curTail.next.get()
            // Хвост отстаёт – продвинем
//            if (tailNext != null) {
//                tail.compareAndSet(curTail, tailNext)
//                continue
//            }
            newNode.prev.set(curTail)
            if (curTail.next.compareAndSet(null, newNode)) {
                tail.compareAndSet(curTail, newNode)
                if (curTail.extractedOrRemoved) curTail.remove()
                return
            }
            else {
                tail.compareAndSet(curTail, curTail.next.get())
                if (curTail.extractedOrRemoved) curTail.remove()
            }
//            if (curTail != tail.get()) {
//                continue
//            }
//            if (curTail.next.compareAndSet(null, newNode)) {
//                // Если curTail был логически удалён, теперь он уже не хвост (у него появился next),
//                // и мы можем попытаться физически его вырезать.
//                if (curTail.extractedOrRemoved && (curTail.element != null && curTail.prev.get() != null)) {
//                    // best-effort: в remove() вторая фаза вырежет curTail,
//                    // если он действительно уже не tail.
//                    curTail.helpRemove()
//                }
//                val tt = tail.get().next.get()
//                if (tt != null) {
//                    if (!tail.compareAndSet(curTail, tt)) {
//                        while (true) {
//                            val tt2 = tail.get()
//                            val n = tt2.moveTailActual() // не сработает, если tail будет на dummy
//                            if (n == null) {
//                                break
//                            }
//                            if (tt2 != n) {
//                                if (tail.compareAndSet(tt2, n)) {
//                                    break
//                                }
//                            }
//                        }
//                        // cas() нельзя, так как curTail не тот
//                        // set() поставит на живого, но надо чтобы его никто потом неправильно не опдвинул
//                    } // если не получится, то можно moveToActual(), потому что remove() мог подвинуть tail вперёд
//                    // но если moveToActual() не найдёт живого, то подвинуть вперёд тоже надо
//                }
////                val n = curTail.moveTailActual()
////                if (n != null && tail.get() != n) {
////                    tail.compareAndSet(curTail, n)
////                }
//                return
//            } else {
//                continue
//            }
        }
    }

    override fun dequeue(): E? {
        // TODO: After moving the `head` pointer forward,
        // TODO: mark the node that contains the extracting
        // TODO: element as "extracted or removed", restarting
        // TODO: the operation if this node has already been removed.
        while (true) {
//            val curTail = tail.get()
//            val n = curTail.moveTailActual()
//            if (n != null && tail.get() != n) {
//                tail.compareAndSet(curTail, n)
//            }
            val curHead = head.get()
            val newHead = curHead.next.get() ?: return null
//            if (newHead == null) {
//                curHead.moveTailActual()
//                return null
//            }
            if (head.compareAndSet(curHead, newHead)) {
                newHead.prev.set(null)
                if (newHead.markExtractedOrRemoved()) {
//                    curHead.moveTailActual()
                    return newHead.element
                }
            }
        }
    }

    override fun remove(element: E): Boolean {
        // Traverse the linked list, searching the specified
        // element. Try to remove the corresponding node if found.
        // DO NOT CHANGE THIS CODE.
        var node = head.get()
        while (true) {
            val next = node.next.get()
            if (next == null) return false
            node = next
            if (node.element == element && node.remove()) return true
        }
    }

    /**
     * This is an internal function for tests.
     * DO NOT CHANGE THIS CODE.
     */
    override fun validate() {
        check(head.get().prev.get() == null) {
            "`head.prev` must be null"
        }
        check(tail.get().next.get() == null) {
            "tail.next must be null"
        }
        // Traverse the linked list
        var node = head.get()
        while (true) {
            if (node !== head.get() && node !== tail.get()) {
                check(!node.extractedOrRemoved) {
                    "Removed node with element ${node.element} found in the middle of the queue"
                }
            }
            val nodeNext = node.next.get()
            // Is this the end of the linked list?
            if (nodeNext == null) break
            // Is next.prev points to the current node?
            val nodeNextPrev = nodeNext.prev.get()
            check(nodeNextPrev != null) {
                "The `prev` pointer of node with element ${nodeNext.element} is `null`, while the node is in the middle of the queue"
            }
            check(nodeNextPrev == node) {
                "node.next.prev != node; `node` contains ${node.element}, `node.next` contains ${nodeNext.element}"
            }
            // Process the next node.
            node = nodeNext
        }
    }

    private class Node<E>(
        var element: E?,
        prev: Node<E>?,
        private val tailRef: AtomicReference<Node<E>>? = null, // AtomicReference<Node<E>>? = null
        private val headRef: AtomicReference<Node<E>>? = null // AtomicReference<Node<E>>? = null
    ) {
        val next = AtomicReference<Node<E>?>(null)
        val prev = AtomicReference(prev)

        /**
         * TODO: Both [dequeue] and [remove] should mark
         * TODO: nodes as "extracted or removed".
         */
        private val _extractedOrRemoved = AtomicBoolean(false)
        val extractedOrRemoved
            get() =
                _extractedOrRemoved.get()

        fun markExtractedOrRemoved(): Boolean =
            _extractedOrRemoved.compareAndSet(false, true)

        fun unmarkExtractedOrRemoved(): Boolean =
            _extractedOrRemoved.compareAndSet(true, false)

//        fun updateLincks(t)

        /**
         * Removes this node from the queue structure.
         * Returns `true` if this node was successfully
         * removed, or `false` if it has already been
         * removed by [remove] or extracted by [dequeue].
         */
        fun remove(): Boolean {
            val removed = markExtractedOrRemoved()
            val curNext = next.get()?: return removed
            val curPrev = prev.get()?: return removed
            curPrev.next.compareAndSet(this, curNext)
            curNext.prev.compareAndSet(this, curPrev)

            if (curPrev.extractedOrRemoved){
                curPrev.remove()
            }
            if (curNext.extractedOrRemoved){
                curNext.remove()
            }
            return removed
        }

//        fun remove(): Boolean {
//            if (!markExtractedOrRemoved()) {// не очень понимаю, зачем мне эта строка
//                return false
//            }
//            while (true) {
//                val prevNode = prev.get()
//                val nextNode = next.get()
//                if (prevNode == null) { //never happen (only for dummy)
//                    // need to move head -> it is not allowed for node
//                    return true
//                }
//                if (nextNode == null) {
//                    prevNode.next.compareAndSet(this, null)
//                    if (tailRef?.get() != headRef?.get()) { // повторить надо ли?
//                        tailRef?.compareAndSet(this, prevNode)
//                    }
//                    moveTailActual()
//                    return true
//                }
//                if (prevNode.next.compareAndSet(this, nextNode)) {
//                    nextNode.prev.compareAndSet(this, prevNode)
//                }
//                moveTailActual()
//                // если head на мне и я его удаляю и есть следующий, то надо head подвинуть
//                val ch = headRef?.get()
//                val chNext = ch?.next?.get()
//                if (ch == this && chNext != null) {
//                    headRef.compareAndSet(ch, chNext)
//                    headRef.get().prev.set(null)
//                }
//                return true
//            }
//        }

        fun helpRemove() {
            val prevNode = prev.get()
            val nextNode = next.get()

            if (prevNode == null) {
                // it is dummy node
                return
            }
            prevNode.next.set(nextNode) //compareAndSet(this, nextNode)
            if (nextNode != null) {
                nextNode.prev.compareAndSet(this, prevNode)
            }
        }

        fun findLiveNodeLast(): Node<E>? {
            var checkNode: Node<E>? = this.next.get()
            var last: Node<E>? = null
            while (checkNode != null) {
                if (!checkNode.extractedOrRemoved) {
                    last = checkNode
                }
                checkNode = checkNode.next.get()
            }
            return last
        }

        fun moveTailActual(): Node<E>? {
            while (true) {
                val curTail = tailRef?.get()
                val n = findLiveNodeLast()
                if (n == null) {
                    moveTail2Head()
                    return n
                }
                if (tailRef == null) {
                    return n
                }
                if (tailRef.compareAndSet(curTail, n)) {
                    return n
                }
            }
        }

        fun moveTail2Head() {
            while (true) {
                var f: Boolean = false
                val start = tailRef?.get()
                var check: Node<E>? = start
                while (check != null) {
                    if (check == headRef?.get()) {
                        if (tailRef?.compareAndSet(start, check) == true) {
                            return
                        } else {
                            f = true
                        }
                    }
                    check = check.next.get()
                }
                if (f) {
                    continue
                }
                return
            }
        }

        fun moveTailRemove() {

        }
        fun isDummy(): Boolean {
            return element == null || prev.get() == null
        }
    }
}

/*
head ->
 */
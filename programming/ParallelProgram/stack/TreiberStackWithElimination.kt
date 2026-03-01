import java.util.concurrent.atomic.*
import kotlin.random.Random

/**
 * @author Кубеш Сергей
 */
open class TreiberStackWithElimination<E> : Stack<E> {
    private val stack = TreiberStack<E>()

    // TODO: Try to optimize concurrent push and pop operations,
    // TODO: synchronizing them in an `rendezvousSlots` cell.
    private val rendezvousSlots = AtomicReferenceArray<Any?>(ELIMINATION_ARRAY_SIZE)

    override fun push(element: E) {
        if (tryPushWithElimination(element)) return
        stack.push(element)
    }

    protected open fun tryPushWithElimination(element: E): Boolean {
        // TODO: Choose a random cell in `rendezvousSlots`
        // TODO: and try to install the element there.
        // TODO: Wait `ELIMINATION_WAIT_CYCLES` loop cycles
        // TODO: in hope that a concurrent `pop()` grabs the
        // TODO: element. If so, clean the cell and finish,
        // TODO: returning `true`. Otherwise, move the cell
        // TODO: to the empty state and return `false`.

        var isInSlot = false
        var slotIndex = 0
        for (i in 0 until ELIMINATION_ARRAY_SIZE) {
            if (rendezvousSlots.compareAndSet(i, null, element)) {
                isInSlot = true
                slotIndex = i
                break
            }
        }
        if (!isInSlot) {
            return false
        }
        for (i in 0 until ELIMINATION_WAIT_CYCLES) {
            if (rendezvousSlots[slotIndex] == null) {
                return true
            }
        }
        return !rendezvousSlots.compareAndSet(slotIndex, element, null)
    }

    override fun pop(): E? = tryPopWithElimination() ?: stack.pop()

    private fun tryPopWithElimination(): E? {
        // TODO: Choose a random cell in `rendezvousSlots`
        // TODO: and try to retrieve an element from there.
        // TODO: On success, return the element.
        // TODO: Otherwise, if the cell is empty, return `null`.
        for (i in 0 until ELIMINATION_ARRAY_SIZE) {
            val elem = rendezvousSlots.get(i)
            if (elem != null && rendezvousSlots.compareAndSet(i, elem, null)) {
                return elem as E?
            }
//            if (rendezvousSlots[i] != null) {
//                val res = rendezvousSlots[i] as E?
////                if (rendezvousSlots[i])
//                rendezvousSlots[i] = null
//                return res
//            }
        }
        return null
    }

    companion object {
        private const val ELIMINATION_ARRAY_SIZE = 3 // Do not change!
        private const val ELIMINATION_WAIT_CYCLES = 1 // Do not change!

        // Initially, all cells are in EMPTY state.
        private val CELL_STATE_EMPTY = null

        // `tryPopElimination()` moves the cell state
        // to `RETRIEVED` if the cell contains an element.
        private val CELL_STATE_RETRIEVED = Any()
    }
}

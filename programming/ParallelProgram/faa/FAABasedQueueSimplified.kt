import java.util.concurrent.atomic.*
import kotlin.math.*

/**
 * @author Кубеш Сергей
 */
class FAABasedQueueSimplified<E> : Queue<E> {
    private val infiniteArray = AtomicReferenceArray<Any?>(1024) // conceptually infinite array
    private val enqIdx = AtomicLong(0)
    private val deqIdx = AtomicLong(0)

    override fun enqueue(element: E) {
        // TODO: Increment the counter atomically via Fetch-and-Add.
        // TODO: Use `getAndIncrement()` function for that.
//        val i = enqIdx.getAndIncrement()
//        val elem = infiniteArray.get(i.toInt())
//        enqIdx.set(i + 1)
        // TODO: Atomically install the element into the cell
        // TODO: if the cell is not poisoned.
//        infiniteArray.compareAndSet(i.toInt(), elem, element)
//        infiniteArray.set(i.toInt(), element
    //        )
        while (true) {
            val i = enqIdx.getAndIncrement()
            if (infiniteArray.compareAndSet(i.toInt(), null, element)) {
                return
            }
        }
    }

    @Suppress("UNCHECKED_CAST")
    override fun dequeue(): E? {
        while (true) {
                // Is this queue empty?
            if (deqIdx.get() >= enqIdx.get()) return null
//            if (enqIdx.get() <= deqIdx.get()) return null
            // TODO: Increment the counter atomically via Fetch-and-Add.
            // TODO: Use `getAndIncrement()` function for that.
            val i = deqIdx.getAndIncrement().toInt()
            // TODO: Try to retrieve an element if the cell contains an
            // TODO: element, poisoning the cell if it is empty.




//            val currentDeqIdx = deqIdx.get()
//            val element = infiniteArray.get(currentDeqIdx.toInt()) //?: continue
//            if (element == null) {
//                if (enqIdx.get() <= currentDeqIdx) {
//                    return null
//                }
//                continue
//            }
//            if (element == POISONED) {
//                deqIdx.compareAndSet(currentDeqIdx, currentDeqIdx + 1)
//                continue
//            }
//            if (deqIdx.compareAndSet(currentDeqIdx, currentDeqIdx + 1)) {
//                if (infiniteArray.compareAndSet(currentDeqIdx.toInt(), element, POISONED)) {
//                    return element as E
//                }
//            }

//            val element = infiniteArray.get(i)
//            if (infiniteArray.compareAndSet(i, null, POISONED)) {
//                continue
//            }
//            if (infiniteArray.compareAndSet(i, element, null)) {
//                return element as E
//            }


//            if (element != POISONED) {//element != null &&
            if (infiniteArray.compareAndSet(i, null, POISONED)) {
                continue
            }
            val element = infiniteArray.get(i)
            if (infiniteArray.compareAndSet(i, element, POISONED)) {
                return element as E
            }
//            }
        }
    }

    override fun validate() {
        for (i in 0 until min(deqIdx.get().toInt(), enqIdx.get().toInt())) {
            check(infiniteArray[i] == null || infiniteArray[i] == POISONED) {
                "`infiniteArray[$i]` must be `null` or `POISONED` with `deqIdx = ${deqIdx.get()}` at the end of the execution"
            }
        }
        for (i in max(deqIdx.get().toInt(), enqIdx.get().toInt()) until infiniteArray.length()) {
            check(infiniteArray[i] == null || infiniteArray[i] == POISONED) {
                "`infiniteArray[$i]` must be `null` or `POISONED` with `enqIdx = ${enqIdx.get()}` at the end of the execution"
            }
        }
    }
}

// TODO: poison cells with this value.
private val POISONED = Any()

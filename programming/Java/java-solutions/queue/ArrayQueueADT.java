package queue;

import java.util.Objects;
import java.util.function.Predicate;

public class ArrayQueueADT {
    private Object[] elements = new Object[2];
    private int size = 2;
    private int tail = -1;
    private int head = 0;

    //Pred: true
    //Post: R == queue (Object[2]) && n' = 2 && n[0] == n[1] == null
    public static ArrayQueueADT create() {
        final ArrayQueueADT queue  = new ArrayQueueADT();
        queue.elements = new Object[2];
        return queue;
    }

    //Pred: x != null
    //Post: n' = n' + 1 && a[n'] == x && immutable(x)
    public static void enqueue(final ArrayQueueADT queue, final Object x) {
        Objects.requireNonNull(x);
        fullQueue(queue);
        queue.tail = (queue.tail + 1) % queue.size;
        queue.elements[(queue.tail)] = x;
        //toStr(queue);
    }

    private static void fullQueue(ArrayQueueADT queue) {
        if (((queue.head - queue.tail == 1) || (queue.head == 0 && queue.tail == queue.size - 1)) && (queue.elements[queue.head] != null)) {
            Object[] elementM = new Object[queue.size * 2];
            System.arraycopy(queue.elements, queue.head, elementM, 0, queue.size - queue.head);
            if (queue.tail < queue.head) {
                System.arraycopy(queue.elements, 0, elementM, queue.size - queue.head, queue.tail + 1);
            }
            queue.head = 0;
            queue.tail = queue.size-1;
            queue.size *= 2;
            queue.elements = elementM;
        }
    }

    //Pred: n >= 1 && queue != null
    //Post: R == a[n] && immutable(n) && n' = n
    public static Object element(ArrayQueueADT queue) {
        Objects.requireNonNull(queue.elements[queue.head]);
        return queue.elements[queue.head];
    }

    //Pred: n >= 1
    //Post: n' == n - 1 && immutable(n') && R = a[n]
    public static Object dequeue(ArrayQueueADT queue) {
        Objects.requireNonNull(queue.elements[queue.head]);
        Object rezult = queue.elements[queue.head];
        queue.elements[queue.head] = null;
        queue.head = (queue.head + 1) % queue.size;
        return rezult;
    }

    //Pred: true
    //Post: R == n && n' == n && immutable(n)
    public static int size(ArrayQueueADT queue) {
        if (queue.tail < queue.head && queue.elements[queue.head] != null) {
            return queue.size - queue.head + queue.tail + 1;
        } else {
            if (queue.elements[queue.head] == null) {
                return 0;
            }
            return (queue.tail - queue.head) % queue.size + 1;
        }
    }

    //Pred: true
    //Post: R == (n == 0) && n' == n && immutable(n)
    public static boolean isEmpty(ArrayQueueADT queue) {
        return (queue.head == (queue.tail + 1) % queue.size) && (queue.elements[queue.head] == null);
    }

    //Pred: queue != null
    //Post: n' == 0
    public static void clear(ArrayQueueADT queue) {
        if (queue.head < queue.tail) {
            for (int i = queue.head; i <= queue.tail; i++) {
                queue.elements[i] = null;
            }
        } else {
            for (int i = queue.head; i != (queue.tail + 1) % queue.size; i = (i + 1) % queue.size) {
                queue.elements[i] = null;
            }
        }
        queue.head = 0;
        queue.tail = -1;
    }

    //Pred: queue != null
    //Post: (-1 <= R <= size(queue)) && n' == n && immutable(n)
    public static int indexIf(ArrayQueueADT queue, final Predicate<Object> predicate) {
        if (queue.head < queue.tail) {
            for (int i = queue.head; i <= queue.tail; i++) {
                if (predicate.test(queue.elements[i])) {
                    return i - queue.head;
                }
            }
        } else {
            int index = 0;
            for (int i = queue.head; i != queue.tail % queue.size; i = (i + 1) % queue.size) {
                if (predicate.test(queue.elements[i])) {
                    return index;
                }
                index++;
            }if (predicate.test(queue.elements[queue.tail])) {
                return index;
            }
        }
        return -1;
    }

    //Pred: queue != null
    //Post: (-1 <= R <= size(queue)) && n' == n && immutable(n)
    public static int lastIndexIf(ArrayQueueADT queue, final Predicate<Object> predicate) {
        if (queue.head < queue.tail) {
            for (int i = queue.tail; i >= queue.head; i--) {
                if (predicate.test(queue.elements[i])) {
                    return i - queue.head;
                }
            }
        } else {
            int index = size(queue) - 1;
            for (int i = queue.tail; i >= 0; i--) {
                if (predicate.test(queue.elements[i])) {
                    return index;
                }
                index--;
            }
            for (int i = queue.size - 1; i >= queue.head; i--) {
                if (predicate.test(queue.elements[i])) {
                    return index;
                }
                index--;
            }
        }
        return -1;
    }
}

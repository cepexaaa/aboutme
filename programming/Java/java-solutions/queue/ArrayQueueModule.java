package queue;

import java.util.Objects;
import java.util.function.Predicate;

public class ArrayQueueModule {
    private static Object[] elements = new Object[20];
    private static int size = 20;
    private static int tail = -1;
    private static int head = 0;

    //Pred: x != null
    //Post: n' = n' + 1 && a[n'] == x && immutable(x)
    public static void enqueue(final Object x) {
        Objects.requireNonNull(x);
        fullQueue();
        tail = (tail + 1) % size;
        elements[tail] = x;
    }

    private static void fullQueue() {
        if (((head - tail == 1) || (head == 0 && tail == size - 1)) && (elements[head] != null)) {//
            Object[] elementM = new Object[size * 2];
            System.arraycopy(elements, head, elementM, 0, size - head);
            if (tail < head) {
                System.arraycopy(elements, 0, elementM, size - head, tail + 1);
            }
            head = 0;
            tail = size-1;
            size*=2;
            elements = elementM;
        }
    }

    //Pred: n >= 1 && queue != null
    //Post: R == a[n] && immutable(n) && n' = n
    public static Object element() {
        Objects.requireNonNull(elements[head]);
        return elements[head];
    }

    //Pred: n >= 1
    //Post: n' == n - 1 && immutable(n') && R = a[n]
    public static Object dequeue() {
        Objects.requireNonNull(elements[head]);
        Object rezult = elements[head];
        elements[head] = null;
        head = (head + 1) % size;
        return rezult;
    }

    //Pred: true
    //Post: R == n && n' == n && immutable(n)
    public static int size() {
        if (tail < head && elements[head] != null) {
            return size - head + tail + 1;
        } else {
            if (elements[head] == null) {
                return 0;
            }
            return (tail - head) % size + 1;
        }
    }

    //Pred: true
    //Post: R == (n == 0) && n' == n && immutable(n)
    public static boolean isEmpty() {
        return (head == (tail + 1) % size) && (elements[head] == null);
    }

    //Pred: queue != null
    //Post: n' == 0
    public static void clear() {
        if (head < tail) {
            for (int i = head; i <= tail; i++) {
                elements[i] = null;
            }
        } else {
            for (int i = head; i != (tail + 1) % size; i = (i + 1) % size) {
                elements[i] = null;
            }
        }
        head = 0;
        tail = -1;
    }

    //Pred: queue != null
    //Post: (-1 <= R <= size(queue)) && n' == n && immutable(n)
    public static int indexIf(final Predicate<Object> predicate) {
        if (head < tail) {
            for (int i = head; i <= tail; i++) {
                if (predicate.test(elements[i])) {
                    return i - head;
                }
            }
        } else {
            int index = 0;
            for (int i = head; i != tail% size; i = (i + 1) % size) {
                if (predicate.test(elements[i])) {
                    return index;
                }
                index++;
            }if (predicate.test(elements[tail])) {
                return index;
            }
        }
        return -1;
    }

    //Pred: queue != null
    //Post: (-1 <= R <= size(queue)) && n' == n && immutable(n)
    public static int lastIndexIf(final Predicate<Object> predicate) {
        int index = size() - 1;
        if (head < tail) {
            for (int i = tail; i >= head; i--) {
                if (predicate.test(elements[i])) {
                    return i - head;
                }
            }
        } else {
            for (int i = tail; i >= 0; i--) {
                if (predicate.test(elements[i])) {
                    return index;
                }
                index--;
            }
            for (int i = size - 1; i >= head; i--) {
                if (predicate.test(elements[i])) {
                    return index;
                }
                index--;
            }
        }
        return -1;
    }
}

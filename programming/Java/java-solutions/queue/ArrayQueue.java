package queue;

import java.util.Objects;
import java.util.Set;
import java.util.function.Predicate;
import java.util.HashSet;

public class ArrayQueue extends AbstractQueue {
    private Object[] elements = new Object[2];
    private int tail = -1;
    private int head = 0;

    //Pred: true
    //Post: R == queue (Object[2]) && n' = 2 && n[0] == n[1] == null
    public static ArrayQueue create() {
        final ArrayQueue queue  = new ArrayQueue();
        queue.elements = new Object[2];
        return queue;
    }

    // Pre: element != null
    // Post: n' = n + 1 && a'[n'] = element && immutable(n)
    protected void enqueueOwn(Object x) {
        fullQueue();
        tail = (tail + 1) % elements.length;
        elements[tail] = x;
    }

    private void fullQueue() {
        if (((head - tail == 1) || (head == 0 && tail == elements.length - 1)) && (elements[head] != null)) {
            Object[] elementM = new Object[elements.length * 2];
            System.arraycopy(elements, head, elementM, 0, elements.length - head);
            if (tail < head) {
                System.arraycopy(elements, 0, elementM, elements.length - head, tail + 1);
            }
            head = 0;
            tail = elements.length-1;
            elements = elementM;
        }
    }

    //Pred: n >= 1 && queue != null
    //Post: R == a[0] && immutable(n) && n' = n
    public Object elementOwn() {
        Objects.requireNonNull(elements[head]);
        return elements[head];
    }

    //Pred: n >= 1
    //Post: n' == n - 1 && immutable(n') && R = a[0]
    public Object dequeueOwn() {
        Objects.requireNonNull(elements[head]);
        Object result = elements[head];
        elements[head] = null;
        head = (head + 1) % elements.length;
        return result;
    }

    //Pred: queue != null
    //Post: n' == 0
    public void clearOwn() {
        if (head < tail) {
            for (int i = head; i <= tail; i++) {
                elements[i] = null;
            }
        } else {
            for (int i = head; i != (tail + 1) % elements.length; i = (i + 1) % elements.length) {
                elements[i] = null;
            }
        }
        head = 0;
        tail = -1;
    }

    //Pred: true
    //Post: (-1 <= R <= size(queue)) && n' == n && immutable(n)
    //R = min(i) : ({predicate(queue[i]) == true})
    public int indexIf(final Predicate<Object> predicate) {
        if (head < tail) {
            for (int i = head; i <= tail; i++) {
                if (predicate.test(elements[i])) {
                    return i - head;
                }
            }
        } else {
            int index = 0;
            for (int i = head; i != tail % elements.length; i = (i + 1) % elements.length) {
                if (predicate.test(elements[i])) {
                    return index;
                }
                index++;
            } if (predicate.test(elements[tail])) {
                return index;
            }
        }
        return -1;
    }

    //Pred: queue != null
    //Post: (-1 <= R <= size(queue)) && n' == n && immutable(n)
    //R = max(i) : ({predicate(queue[i]) == true})
    public int lastIndexIf(final Predicate<Object> predicate) {
        int index = size - 1;
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
            for (int i = elements.length - 1; i >= head; i--) {
                if (predicate.test(elements[i])) {
                    return index;
                }
                index--;
            }
        }
        return -1;
    }

    public void distinctOwn() {
        Set<Object> set = new HashSet<>();
        Object[] uniqueElements = new Object[elements.length];
        int uniqueIndex = 0;
        if (head < tail) {
            for (int i = head; i <= tail; i++) {
                if (!set.contains(elements[i]) && elements[i] != null) {
                    uniqueElements[uniqueIndex++] = elements[i];
                }set.add(elements[i]);
            }
        } else {
            for (int i = head; i != tail % elements.length; i = (i + 1) % elements.length) {
                if (!set.contains(elements[i]) && elements[i] != null) {
                    uniqueElements[uniqueIndex++] = elements[i];
                }set.add(elements[i]);
            }if (!set.contains(elements[tail]) && elements[tail] != null) {
                uniqueElements[uniqueIndex++] = elements[tail];
            }
        }
        head = 0;
        tail = uniqueIndex-1;
        elements = uniqueElements;
        size = uniqueIndex;
    }
}

package queue;

import java.util.HashSet;
import java.util.Objects;
import java.util.Set;


public class LinkedQueue extends AbstractQueue {
    private Node head;
    private Node tail;

    // Pre: element != null
    // Post: n' = n + 1 && a'[n'] = element && immutable(n)
    protected void enqueueOwn(Object x) {
        Node newNode = new Node(x, null);
        if (tail == null) {
            head = newNode;
        } else {
            tail.next = newNode;
        }
        tail = newNode;
    }

    // Pre: n > 0
    // Post: R = a[1] && n' = n - 1 && immutable(n')
    public Object dequeueOwn() {
        Objects.requireNonNull(head);
        Object result = head.value;
        head = head.next;
        if (head == null) {
            tail = null;
        }
        return result;
    }

    // Pre: n > 0
    // Post: R = a[1] && n' = n && immutable(n)
    public Object elementOwn() {
        Objects.requireNonNull(head);
        return head.value;
    }

    //Pred: queue != null
    //Post: n' == 0
    public void clearOwn() {
        head = null;
        tail = null;
    }

    private static class Node {
        private final Object value;
        private Node next;

        public Node(Object value, Node next) {
            assert value != null;

            this.value = value;
            this.next = next;
        }
    }

    public void distinctOwn() {
        Set<Object> set = new HashSet<>();
        LinkedQueue uniqueElements = new LinkedQueue();
        Node current = head;
        while (current != null) {
            if (!set.contains(current.value)) {
                uniqueElements.enqueue(current.value);
                set.add(current.value);
            }
            current = current.next;
        }
        head = uniqueElements.head;
        tail = uniqueElements.tail;
        size = uniqueElements.size;
    }
}

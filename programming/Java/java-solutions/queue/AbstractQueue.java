package queue;

import java.util.Objects;

public abstract class AbstractQueue implements Queue{
    protected int size;

    public void enqueue(final Object x) {
        Objects.requireNonNull(x);
        size++;
        enqueueOwn(x);
    }
    protected abstract void enqueueOwn(Object x);

    public boolean isEmpty() {
        return (0 == size());
    }
    public Object dequeue() {
        size--;
         return dequeueOwn();
    }
    protected abstract Object dequeueOwn();
    public Object element() {
        return elementOwn();
    }

    protected abstract Object elementOwn();
    public int size() {
        return size;
    };
    public void clear() {
        size = 0;
        clearOwn();
    }
    protected abstract void clearOwn();
    public void distinct() {
        distinctOwn();
    }

    protected abstract void distinctOwn();

}

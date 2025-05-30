package info.kgeorgiy.ja.kubesh.arrayset;

import java.util.*;

public class ArraySet<E> extends AbstractSet<E> implements SortedSet<E> {

    private final List<E> elements;
    private final Comparator<? super E> comparator;

    public ArraySet() {
        this(Collections.emptyList(), null);
    }

    public ArraySet(Collection<? extends E> collection) {
        this(collection, null);
    }

    public ArraySet(Comparator<? super E> comparator) {
        this(Collections.emptyList(), comparator);
    }

    public ArraySet(Collection<? extends E> collection, Comparator<? super E> comparator) {
        TreeSet<E> treeSet = new TreeSet<>(comparator);
        treeSet.addAll(collection);
        this.elements = List.copyOf(treeSet);
        this.comparator = comparator;
    }

    private ArraySet(List<E> elements, Comparator<? super E> comparator) {
        this.elements = elements;
        this.comparator = comparator;
    }


    @Override
    public Comparator<? super E> comparator() {
        return comparator;
    }

    @Override
    public SortedSet<E> subSet(E fromElement, E toElement) {
        if (fromElement == null || toElement == null) {
            throw new NullPointerException("Elements cannot be null");
        }
        if (checkArguments(fromElement, toElement, false)) {
            throw new IllegalArgumentException("fromElement > toElement");
        }
        int fromIndex = findIndex(fromElement);
        int toIndex = findIndex(toElement);
        if (fromIndex > toIndex) {
            throw new IllegalArgumentException("fromElement > toElement");
        }
        return new ArraySet<>(elements.subList(fromIndex, toIndex), comparator);
    }

    private SortedSet<E> subSetImpl(E toElement, boolean isHead) {
        if (toElement == null) {
            throw new NullPointerException("Elements cannot be null");
        }
        int toIndex = findIndex(toElement);
        if (isHead) {
            return new ArraySet<>(elements.subList(0, toIndex), comparator);
        } else {
            return new ArraySet<>(elements.subList(toIndex, elements.size()), comparator);
        }

    }

    @Override
    public SortedSet<E> headSet(E toElement) {
        return subSetImpl(toElement, true);
    }

    @Override
    public SortedSet<E> tailSet(E fromElement) {
        return subSetImpl(fromElement, false);
    }

    @Override
    public E first() {
        if (elements.isEmpty()) {
            throw new NoSuchElementException();
        }
        return elements.getFirst();
    }

    @Override
    public E last() {
        if (elements.isEmpty()) {
            throw new NoSuchElementException();
        }
        return elements.getLast();
    }


    @Override
    public Iterator<E> iterator() {
        return elements.iterator();
    }

    @Override
    public int size() {
        return elements.size();
    }

    @Override
    @SuppressWarnings("unchecked")
    public boolean contains(Object o) {
        return Collections.binarySearch(elements, (E) Objects.requireNonNull(o), comparator) >= 0;
    }


    private int findIndex(E element) {
        int index = Collections.binarySearch(elements, Objects.requireNonNull(element), comparator);
        if (index < 0) {
            index = -index - 1;
        }
        return index;
    }

    private boolean checkArguments(E fromElement, E toElement, boolean edgeSubSet) {
        if (!edgeSubSet){
            if (comparator != null && comparator.compare(fromElement, toElement) > 0) {
                return true;
            } else if (comparator == null) {
                return Collections.reverseOrder().reversed().compare(fromElement, toElement) > 0;
            }
        }
        return false;
    }
}


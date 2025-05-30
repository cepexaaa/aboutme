package info.kgeorgiy.ja.kubesh.iterative;

import info.kgeorgiy.java.advanced.iterative.ScalarIP;
import info.kgeorgiy.java.advanced.mapper.ParallelMapper;

import java.util.*;
import java.util.function.Function;
import java.util.function.Predicate;

/**
 * Scalar iterative parallelism support.
 * Implements ScalarIP.
 * {@inheritDoc}
 *
 * @see info.kgeorgiy.java.advanced.iterative.ScalarIP
 */
public class IterativeParallelism implements ScalarIP {
    private final ParallelMapper parallelMapper;

    /**
     * Constructor is used to get ParallelMapper.
     * Input parallelMapper allow to use map, which apply function to list parallel.
     *
     * @param mapper class with map, which can use parallel execution.
     * @see ParallelMapper
     */
    public IterativeParallelism(ParallelMapper mapper) {
        parallelMapper = mapper;
    }

    /**
     * Default constructor
     */
    public IterativeParallelism() {
        parallelMapper = null;
    }

    private <T, R> List<R> parallelExecute(int threads, List<T> values, Function<Range, R> task) throws InterruptedException {
        threads = Math.min(threads, values.size());
        Range[] ranges = splitRange(values.size(), threads);

        if (parallelMapper != null) {
            return parallelMapper.map(task, List.of(ranges));
        } else {
            List<R> results = new ArrayList<>(Collections.nCopies(threads, null));
            Thread[] workers = new Thread[threads];

            for (int i = 0; i < threads; i++) {
                final int threadId = i;
                workers[i] = new Thread(() -> results.set(threadId, task.apply(ranges[threadId])));
                workers[i].start();
            }
            ParallelMapperImpl.joinAll(List.of(workers));
            return results;
        }
    }

    private <T> int findArg(int threads, List<T> values, Comparator<? super T> comparator, boolean findMax)
            throws InterruptedException {
        if (values.isEmpty()) {
            throw new NoSuchElementException("List is empty");
        }
        List<ExtremumResult<T>> results = parallelExecute(threads, values, range -> findExtremumInRange(values, range.start, range.end, comparator, findMax));
        ExtremumResult<T> globalResult = findExtremumInList(results, comparator, findMax);
        return globalResult.index;
    }

    private <T> ExtremumResult<T> findExtremumInList(List<ExtremumResult<T>> results, Comparator<? super T> comparator, boolean findMax) {
        ExtremumResult<T> globalResult = results.getFirst();
        for (int i = 1; i < results.size(); i++) {
            int cmp = comparator.compare(results.get(i).value, globalResult.value);
            if ((findMax && cmp > 0) || (!findMax && cmp < 0)) {
                globalResult = results.get(i);
            }
        }
        return globalResult;
    }

    private <T> ExtremumResult<T> findExtremumInRange(List<T> values, int start, int end, Comparator<? super T> comparator, boolean findMax) {
        T extremum = values.get(start);
        int index = start;
        for (int j = start + 1; j < end; j++) {
            T current = values.get(j);
            int cmp = comparator.compare(current, extremum);
            if ((findMax && cmp > 0) || (!findMax && cmp < 0)) {
                extremum = current;
                index = j;
            }
        }
        return new ExtremumResult<>(extremum, index);
    }

    private <T> int findIndex(int threads, List<T> values, Predicate<? super T> predicate, boolean findLast)
            throws InterruptedException {
        if (values.isEmpty()) {
            return -1;
        }
        List<Integer> results = parallelExecute(threads, values, range -> findIndexInRange(values, range.start, range.end, predicate, findLast));
        return findGlobalIndex(results, findLast);
    }

    private <T> Integer findIndexInRange(List<T> values, int start, int end, Predicate<? super T> predicate, boolean findLast) {
        if (findLast) {
            for (int j = end - 1; j >= start; j--) {
                if (predicate.test(values.get(j))) {
                    return j;
                }
            }
        } else {
            for (int j = start; j < end; j++) {
                if (predicate.test(values.get(j))) {
                    return j;
                }
            }
        }
        return -1;
    }

    private int findGlobalIndex(List<Integer> indices, boolean findLast) {
        int result = -1;
        for (int index : indices) {
            if (index != -1) {
                if (findLast) {
                    result = Math.max(result, index);
                } else if (result == -1 || index < result) {
                    result = index;
                }
            }
        }
        return result;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> int argMax(int threads, List<T> values, Comparator<? super T> comparator) throws InterruptedException {
        return findArg(threads, values, comparator, true);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> int argMin(int threads, List<T> values, Comparator<? super T> comparator) throws InterruptedException {
        return findArg(threads, values, comparator, false);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> int indexOf(int threads, List<T> values, Predicate<? super T> predicate) throws InterruptedException {
        return findIndex(threads, values, predicate, false);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> int lastIndexOf(int threads, List<T> values, Predicate<? super T> predicate) throws InterruptedException {
        return findIndex(threads, values, predicate, true);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> long sumIndices(int threads, List<? extends T> values, Predicate<? super T> predicate)
            throws InterruptedException {
        if (values.isEmpty()) {
            return 0;
        }
        List<Long> partialSums = parallelExecute(
                threads,
                values,
                range -> {
                    long sum = 0;
                    for (int i = range.start; i < range.end; i++) {
                        if (predicate.test(values.get(i))) {
                            sum += i;
                        }
                    }
                    return sum;
                }
        );
        long totalSum = 0;
        for (long sum : partialSums) {
            totalSum += sum;
        }
        return totalSum;
    }

    private static Range[] splitRange(int size, int parts) {
        Range[] ranges = new Range[parts];
        int partSize = size / parts;
        int remainder = size % parts;

        int start = 0;
        for (int i = 0; i < parts; i++) {
            int end = start + partSize + (i < remainder ? 1 : 0);
            ranges[i] = new Range(start, end);
            start = end;
        }
        return ranges;
    }

    private static final class Range {
        final int start;
        final int end;

        Range(int start, int end) {
            this.start = start;
            this.end = end;
        }
    }

    private static class ExtremumResult<T> {
        final T value;
        final int index;

        ExtremumResult(T value, int index) {
            this.value = value;
            this.index = index;
        }
    }
}
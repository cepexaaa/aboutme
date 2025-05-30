package info.kgeorgiy.ja.kubesh.iterative;

import info.kgeorgiy.java.advanced.mapper.ParallelMapper;

import java.util.*;
import java.util.function.Function;

/**
 * Maps function over lists.
 * Implements ParallelMapper
 * {@inheritDoc}
 *
 * @see ParallelMapper
 */
public class ParallelMapperImpl implements ParallelMapper {
    private final TaskQueue taskQueue = new TaskQueue();
    private final List<Thread> workers;

    /**
     * Constructor creates input count threads.
     * This constructor create certain number of threads.
     * Each thread takes the earlier task from queue and complete it.
     *
     * @param threads count of available parallel threads
     */
    public ParallelMapperImpl(int threads) {
        if (threads <= 0) {
            throw new IllegalArgumentException("Number of threads must be positive");
        }

        this.workers = new ArrayList<>(threads);
        for (int i = 0; i < threads; i++) {
            Thread worker = new Thread(() -> {
                try {
                    while (!Thread.interrupted()) {
                        Runnable task = taskQueue.poll();
                        Objects.requireNonNull(task).run();
                    }
                } catch (InterruptedException ignored) {
                    Thread.currentThread().interrupt();
                }
            });
            worker.start();
            workers.add(worker);
        }
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T, R> List<R> map(Function<? super T, ? extends R> f, List<? extends T> items) throws InterruptedException {
        List<Task<R>> tasks = new ArrayList<>(items.size());
        for (int i = 0; i < items.size(); i++) {
            tasks.add(new Task<>());
        }

        for (int i = 0; i < items.size(); i++) {
            final int index = i;
            taskQueue.add(() -> {
                Task<R> task = tasks.get(index);
                try {
                    task.result = f.apply(items.get(index));
                } catch (RuntimeException e) {
                    task.exception = e;
                } finally {
                    task.complete();
                }
            });
        }

        List<R> results = new ArrayList<>(items.size());
        List<Exception> exceptions = new ArrayList<>();
        for (Task<R> task : tasks) {
            task.awaitCompletion();
            results.add(task.result);
            if (task.exception != null) {
                exceptions.add(task.exception);
            }
        }

        if (!exceptions.isEmpty()) {
            Exception primary = exceptions.getFirst();
            for (int i = 1; i < exceptions.size(); i++) {
                primary.addSuppressed(exceptions.get(i));
            }
            System.err.println("Errors during task execution:");
            throw new RuntimeException(primary);
        }

        return results;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void close() {
        workers.forEach(Thread::interrupt);
        try {
            joinAll(workers);
        } catch (InterruptedException e) {
            System.err.println("Interrupted while waiting for worker threads to finish");
            Thread.currentThread().interrupt();
        }
    }

    /**
     * Join all threads.
     * This method try to close each thread.
     * If error is occurred, it will continue join and throw exception in the end of work.
     *
     * @param threads list of worked threads
     * @throws InterruptedException if some thread throw exception
     */
    public static void joinAll(List<Thread> threads) throws InterruptedException {
        InterruptedException exception = null;
        for (Thread thread : threads) {
            while (true) {
                try {
                    thread.join();
                    break;
                } catch (InterruptedException e) {
                    if (exception == null) {
                        exception = new InterruptedException("Interrupted while joining threads");
                    }
                    exception.addSuppressed(e);
                }
            }
        }
        if (exception != null) {
            throw exception;
        }
    }

    private static class TaskQueue {
        private final Queue<Runnable> taskQueue = new LinkedList<>();

        public synchronized Runnable poll() throws InterruptedException {
            while (taskQueue.isEmpty()) {
                wait();
            }
            return taskQueue.poll();
        }

        public synchronized void add(Runnable task) {
            taskQueue.add(task);
            notify();
        }
    }

    private static class Task<R> {
        public R result;
        private boolean completed = false;
        public Exception exception;

        public synchronized void complete() {
                completed = true;
                this.notifyAll();
        }

        public synchronized void awaitCompletion() throws InterruptedException {
                while (!completed) {
                    this.wait();
                }
        }
    }
}


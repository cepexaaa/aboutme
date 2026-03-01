import java.util.concurrent.atomic.AtomicReference;

public class Solution implements Lock<Solution.Node> {
    private final Environment env;
    private final AtomicReference<Node> tail = new AtomicReference<>(null);

    public Solution(Environment env) {
        this.env = env;
    }

    @Override
    public Node lock() {
        Node my = new Node();
        Node prev = tail.getAndSet(my);

        if (prev != null) {
            prev.next.set(my);
            while (my.locked.get()) {
                env.park();
            }
        } else {
            my.locked.set(false);
        }

        return my;
    }

    @Override
    public void unlock(Node node) {
        if (node.next.get() == null) {
            if (tail.compareAndSet(node, null)) {
                return;
            }
            while (node.next.get() == null) {
                // Spin-wait (как и в оригинальном MCS)
            }
        }

        Node next = node.next.get();
        next.locked.set(false);
        env.unpark(next.thread);
    }

    public static class Node {
        final Thread thread = Thread.currentThread(); // Поток, создавший узел
        final AtomicReference<Boolean> locked = new AtomicReference<>(true); // Начинаем заблокированными
        final AtomicReference<Node> next = new AtomicReference<>(null); // Следующий узел в очереди
    }
}
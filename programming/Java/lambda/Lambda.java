package info.kgeorgiy.ja.kubesh.lambda;

import info.kgeorgiy.java.advanced.lambda.EasyLambda;
import info.kgeorgiy.java.advanced.lambda.Trees;

import java.util.*;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;
import java.util.stream.Collector;
import java.util.stream.Collectors;

public class Lambda implements EasyLambda {
//    abstract static class AbstractTreeSpliterator<T, N> implements Spliterator<T> {
//        protected final N tree;
//        protected final Stack<N> stack = new Stack<>();
//
//        protected AbstractTreeSpliterator(N root) {
//            if (root != null) {
//                pushLeft(root);
//            }
//            this.tree = root;
//        }
//
//        protected abstract void pushLeft(N node);
//
//        @Override
//        public boolean tryAdvance(Consumer<? super T> action) {
////            while (!stack.isEmpty()) {
////                N node = stack.pop();
////                if (isTreesLeaf(node)) {
////                    Trees.Leaf<T> leaf = castToLeaf(node);
////                    action.accept(leaf.value());
////                } else {
////                    node = node.getLeft();
////                }
////            }
//            if (stack.isEmpty()) {
//                return false;
//            }
//            N node = stack.pop();
//            if (isTreesLeaf(node)) {
//                Trees.Leaf<T> leaf = castToLeaf(node);
//                action.accept(leaf.value());
//            } else {
//                processNode(node);
//                return tryAdvance(action);
//            }
//            return true;
//        }
//
//        protected abstract boolean isTreesLeaf(N node);
//        protected abstract Trees.Leaf<T> castToLeaf(N node);
//        protected abstract void processNode(N node);
//        protected abstract int getDepth(N node);
//
//        @Override
//        public Spliterator<T> trySplit() {
//            if (stack.isEmpty()) {
//                return null;
//            }
//            return createSplitSpliterator(stack.pop());
//        }
//
//        protected abstract Spliterator<T> createSplitSpliterator(N node);
//
//        @Override
//        public long estimateSize() {
//            if (tree == null) {
//                return -1;
//            }
//            if (isTreesLeaf(tree)) {
//                return 1;
//            } return getDepth(tree);
//        }
//
//        @Override
//        public int characteristics() {
//            return Spliterator.ORDERED | Spliterator.IMMUTABLE;
//        }
//    }
//
//    @Override
//    public <T> Spliterator<T> binaryTreeSpliterator(Trees.Binary<T> tree) {
//        return new AbstractTreeSpliterator<>(tree) {
//            @Override
//            protected void pushLeft(Trees.Binary<T> node) {
//                pushLeftIn(
//                        node,
//                        stack,
//                        n -> n instanceof Trees.Binary.Branch<T>,
//                        n -> ((Trees.Binary.Branch<T>) n).left()
//                );
//            }
//
//            @Override
//            protected void processNode(Trees.Binary<T> node) {
//                if (node instanceof Trees.Binary.Branch<T> branch) {
//                    pushLeft(branch.right());
//                }
//            }
//
//            @Override
//            protected int getDepth(Trees.Binary<T> node) {
//                if (isTreesLeaf(node)) {
//                    return 1;
//                }
//                if (node instanceof Trees.Binary.Branch<T>(Trees.Binary<T> left, Trees.Binary<T> right)) {
//                    int c = 0;
//                    if (right != null) {
//                        if (isTreesLeaf(right)) {
//                            c++;
//                        } else {
//                            return -1;
//                        }
//                    } if (left != null) {
//                        if (isTreesLeaf(left)) {
//                            c++;
//                        } else {
//                            return -1;
//                        }
//                    }
//                    return c;
//                }
//                return -1;
//            }
//
//            @Override
//            protected Spliterator<T> createSplitSpliterator(Trees.Binary<T> node) {
//                return binaryTreeSpliterator(node);
//            }
//
//            protected boolean isTreesLeaf(Trees.Binary<T> node) {
//                return node instanceof Trees.Leaf<T>;
//            }
//
//            @Override
//            protected Trees.Leaf<T> castToLeaf(Trees.Binary<T> node) {
//                return (Trees.Leaf<T>) node;
//            }
//        };
//    }
//
//    @Override
//    public <T> Spliterator<T> sizedBinaryTreeSpliterator(Trees.SizedBinary<T> tree) {
//        return new AbstractTreeSpliterator<>(tree) {
//            private long size = tree.size();
//
//            @Override
//            protected void pushLeft(Trees.SizedBinary<T> node) {
//                pushLeftIn(
//                        node,
//                        stack,
//                        n -> n instanceof Trees.SizedBinary.Branch<T>,
//                        n -> ((Trees.SizedBinary.Branch<T>) n).left()
//                );
//            }
//
//            @Override
//            protected int getDepth(Trees.SizedBinary<T> node) {
//                if (isTreesLeaf(node)) {
//                    return 1;
//                }
//                return -1;
//            }
//
//            @Override
//            protected boolean isTreesLeaf(Trees.SizedBinary<T> node) {
//                return node instanceof Trees.Leaf<T>;
//            }
//
//            @Override
//            protected Trees.Leaf<T> castToLeaf(Trees.SizedBinary<T> node) {
//                return (Trees.Leaf<T>) node;
//            }
//
//            @Override
//            protected void processNode(Trees.SizedBinary<T> node) {
//                if (node instanceof Trees.SizedBinary.Branch<T> branch) {
//                    pushLeft(branch.right());
//                }
//            }
//
//            @Override
//            protected Spliterator<T> createSplitSpliterator(Trees.SizedBinary<T> node) {
//                if (size <= 1) {
//                    return null;
//                }
//                long half = size / 2;
//                size -= half;
//                return sizedBinaryTreeSpliterator(node);
//            }
//
//            @Override
//            public long estimateSize() {
//                return size;
//            }
//
//            @Override
//            public int characteristics() {
//                return Spliterator.ORDERED | Spliterator.SIZED | Spliterator.SUBSIZED | Spliterator.IMMUTABLE;
//            }
//        };
//    }
//
//    @Override
//    public <T> Spliterator<T> naryTreeSpliterator(Trees.Nary<T> tree) {
//        return new AbstractTreeSpliterator<>(tree) {
//            @Override
//            protected void pushLeft(Trees.Nary<T> node) {
//                if (node != null) {
//                    stack.push(node);
//                }
//            }
//
//            @Override
//            protected boolean isTreesLeaf(Trees.Nary<T> node) {
//                return node instanceof Trees.Leaf<T>;
//            }
//
//            @Override
//            protected Trees.Leaf<T> castToLeaf(Trees.Nary<T> node) {
//                return (Trees.Leaf<T>) node;
//            }
//
//            @Override
//            protected void processNode(Trees.Nary<T> node) {
//                if (node instanceof Trees.Nary.Node<T>(List<Trees.Nary<T>> children)) {
//                    for (int i = children.size() - 1; i >= 0; i--) {
//                        stack.push(children.get(i));
//                    }
//                }
//            }
//
//            @Override
//            protected int getDepth(Trees.Nary<T> node) {
//                if (isTreesLeaf(node)) {
//                    return 1;
//                }
//                if (node instanceof Trees.Nary.Node<T>(List<Trees.Nary<T>> branch)) {
//                    int c = 0;
//                    for (Trees.Nary<T> n : branch) {
//                        if (isTreesLeaf(n)) {
//                            c++;
//                        } else {
//                            return -1;
//                        }
//                    }
//                    return c;
//                }
//                return -1;
//            }
//
//            @Override
//            protected Spliterator<T> createSplitSpliterator(Trees.Nary<T> node) {
//                return naryTreeSpliterator(node);
//            }
//        };
//    }

    private <N> void pushLeftIn(
            N node,
            Stack<N> stack,
            Predicate<N> isBranch,
            Function<N, N> getLeft
    ) {
        while (isBranch.test(node)) {
            stack.push(node);
            node = getLeft.apply(node);
        }
        if (node != null) {
            stack.push(node);
        }
    }

    @Override
    public <T> Collector<T, ?, Optional<T>> first() {
        return Collectors.reducing((a, b) -> a);
    }

    @Override
    public <T> Collector<T, ?, Optional<T>> last() {
        return Collectors.reducing((a, b) -> b);
    }

    @Override
    public <T> Collector<T, ?, Optional<T>> middle() {
        return Collector.of(
                ArrayList<T>::new,
                List::add,
                (list1, list2) -> {
                    list1.addAll(list2);
                    return list1;
                },
                list -> {
                    if (list.isEmpty()) {
                        return Optional.empty();
                    }
                    return Optional.of(list.get(list.size() / 2));
                }
        );
    }

    static class SpecialStringBuilder {
        private final StringBuilder stringBuilder;
        private boolean first;

        SpecialStringBuilder() {
            stringBuilder = new StringBuilder();
            first = true;
        }

        public void append(CharSequence c) {
            stringBuilder.append(c);
        }

        public boolean getFlag() {
            return first;
        }

        public void setFlag(boolean flag) {
            first = flag;
        }

        public StringBuilder getStringBuilder() {
            return stringBuilder;
        }

        @Override
        public String toString() {
            return stringBuilder.toString();
        }
    }

    @Override
    public Collector<CharSequence, ?, String> commonPrefix() {
        return getStringCollector(Lambda::findBothPrefLength);
    }

    @Override
    public Collector<CharSequence, ?, String> commonSuffix() {
        return getStringCollector(Lambda::findBothSufLength);
    }

    private static Collector<CharSequence, SpecialStringBuilder, String> getStringCollector(BiConsumer<StringBuilder, CharSequence> lengthFinder) {
        return Collector.of(
                SpecialStringBuilder::new,
                (ssb, cs) -> processSequence(ssb, cs, lengthFinder),
                (ssb1, ssb2) -> {
                    lengthFinder.accept(ssb1.getStringBuilder(), ssb2.getStringBuilder());
                    return ssb1;
                },
                SpecialStringBuilder::toString
        );
    }

    private static void findBothSufLength(StringBuilder sb, CharSequence cs) {
        int length = 0;
        while (sb.length() > length && cs.length() > length && cs.charAt(cs.length() - length - 1) == sb.charAt(sb.length() - length - 1)) {
            length++;
        }
        sb.delete(0, sb.length() - length);
    }

    private static void findBothPrefLength(StringBuilder sb, CharSequence cs) {
        int length = 0;
        while (sb.length() > length && cs.length() > length && cs.charAt(length) == sb.charAt(length)) {
            length++;
        }
        sb.setLength(length);
    }

    private static void processSequence(
            SpecialStringBuilder ssb,
            CharSequence cs,
            BiConsumer<StringBuilder, CharSequence> lengthFinder
    ) {
        if (ssb.getFlag()) {
            ssb.append(cs);
            ssb.setFlag(false);
        } else {
            lengthFinder.accept(ssb.getStringBuilder(), cs);
        }
    }



    //separated classes








    private static class BinaryTreeSpliterator<T> implements Spliterator<T> {
        private final Trees.Binary<T> tree;
        private final Stack<Trees.Binary<T>> stack = new Stack<>();

        BinaryTreeSpliterator(Trees.Binary<T> root) {
            this.tree = root;
            if (root != null) {
                pushLeft(root);
            }
        }

        private void pushLeft(Trees.Binary<T> node) {
            while (node instanceof Trees.Binary.Branch<T> branch) {
                stack.push(node);
                node = branch.left();
            }
            if (node != null) {
                stack.push(node);
            }
        }

        @Override
        public boolean tryAdvance(Consumer<? super T> action) {
            if (stack.isEmpty()) {
                return false;
            }
            Trees.Binary<T> node = stack.pop();
            if (node instanceof Trees.Leaf<T> leaf) {
                action.accept(leaf.value());
            } else if (node instanceof Trees.Binary.Branch<T> branch) {
                pushLeft(branch.right());
                return tryAdvance(action);
            }
            return true;
        }

        @Override
        public Spliterator<T> trySplit() {
            if (stack.isEmpty()) {
                return null;
            }
            return new BinaryTreeSpliterator<>(stack.pop());
        }

        @Override
        public long estimateSize() {
            if (tree == null) {
                return 0;
            }
            if (tree instanceof Trees.Leaf<T>) {
                return 1;
            }
            return getDepth(tree);
        }

        protected int getDepth(Trees.Binary<T> node) {
                if (node instanceof Trees.Leaf<T> leaf) {
                    return 1;
                }
                if (node instanceof Trees.Binary.Branch<T>(Trees.Binary<T> left, Trees.Binary<T> right)) {
                    int c = 0;
                    if (right != null) {
                        if (right instanceof Trees.Leaf<T>) {
                            c++;
                        } else {
                            return -1;
                        }
                    } if (left != null) {
                        if (left instanceof Trees.Leaf<T>) {
                            c++;
                        } else {
                            return -1;
                        }
                    }
                    return c;
                }
                return -1;
            }

        @Override
        public int characteristics() {
            return Spliterator.ORDERED | Spliterator.IMMUTABLE;
        }
    }

    private static class SizedBinaryTreeSpliterator<T> implements Spliterator<T> {
        private final Trees.SizedBinary<T> tree;
        private final Stack<Trees.SizedBinary<T>> stack = new Stack<>();
        private long size;

        SizedBinaryTreeSpliterator(Trees.SizedBinary<T> root) {
            this.tree = root;
            this.size = root != null ? root.size() : 0;
            if (root != null) {
                pushLeft(root);
            }
        }

        private void pushLeft(Trees.SizedBinary<T> node) {
            while (node instanceof Trees.SizedBinary.Branch<T> branch) {
                stack.push(node);
                node = branch.left();
            }
            if (node != null) {
                stack.push(node);
            }
        }

        @Override
        public boolean tryAdvance(Consumer<? super T> action) {
            if (stack.isEmpty()) {
                return false;
            }
            Trees.SizedBinary<T> node = stack.pop();
            if (node instanceof Trees.Leaf<T> leaf) {
                action.accept(leaf.value());
                size--;
            } else if (node instanceof Trees.SizedBinary.Branch<T> branch) {
                pushLeft(branch.right());
                return tryAdvance(action);
            }
            return true;
        }

        @Override
        public Spliterator<T> trySplit() {
            if (size <= 1 || stack.isEmpty()) {
                return null;
            }
            long half = size / 2;
            size -= half;
            return new SizedBinaryTreeSpliterator<>(stack.pop());
        }

        @Override
        public long estimateSize() {
            return size;
        }

        @Override
        public int characteristics() {
            return Spliterator.ORDERED | Spliterator.SIZED | Spliterator.SUBSIZED | Spliterator.IMMUTABLE;
        }
    }

    private static class NaryTreeSpliterator<T> implements Spliterator<T> {
        private final Trees.Nary<T> tree;
        private final Stack<Trees.Nary<T>> stack = new Stack<>();

        NaryTreeSpliterator(Trees.Nary<T> root) {
            this.tree = root;
            if (root != null) {
                stack.push(root);
            }
        }

        @Override
        public boolean tryAdvance(Consumer<? super T> action) {
            if (stack.isEmpty()) {
                return false;
            }
            Trees.Nary<T> node = stack.pop();
            if (node instanceof Trees.Leaf<T> leaf) {
                action.accept(leaf.value());
            } else if (node instanceof Trees.Nary.Node<T> naryNode) {
                List<Trees.Nary<T>> children = naryNode.children();
                for (int i = children.size() - 1; i >= 0; i--) {
                    stack.push(children.get(i));
                }
                return tryAdvance(action);
            }
            return true;
        }

        @Override
        public Spliterator<T> trySplit() {
            if (stack.isEmpty()) {
                return null;
            }
            return new NaryTreeSpliterator<>(stack.pop());
        }

        @Override
        public long estimateSize() {
            if (tree == null) {
                return 0;
            }
            if (tree instanceof Trees.Leaf<T>) {
                return 1;
            }
            return getDepth(tree);
        }

        protected int getDepth(Trees.Nary<T> node) {
                if (node instanceof Trees.Leaf<T>) {
                    return 1;
                }
                if (node instanceof Trees.Nary.Node<T>(List<Trees.Nary<T>> branch)) {
                    int c = 0;
                    for (Trees.Nary<T> n : branch) {
                        if (n instanceof Trees.Leaf<T>) {
                            c++;
                        } else {
                            return -1;
                        }
                    }
                    return c;
                }
                return -1;
            }

        @Override
        public int characteristics() {
            return Spliterator.ORDERED | Spliterator.IMMUTABLE;
        }
    }

    @Override
    public <T> Spliterator<T> binaryTreeSpliterator(Trees.Binary<T> tree) {
        return new BinaryTreeSpliterator<>(tree);
    }

    @Override
    public <T> Spliterator<T> sizedBinaryTreeSpliterator(Trees.SizedBinary<T> tree) {
        return new SizedBinaryTreeSpliterator<>(tree);
    }

    @Override
    public <T> Spliterator<T> naryTreeSpliterator(Trees.Nary<T> tree) {
        return new NaryTreeSpliterator<>(tree);
    }




}



package queue;

import java.util.Arrays;

class ArrayQueueModuleTest {
    public static void main(String[] args) {
        ArrayQueueModule.clear();
        ad(1);
        ArrayQueueModule.clear();
        ad(1);
        ArrayQueueModule.clear();
        ad(1);
        ArrayQueueModule.clear();
        ad(1);
        ArrayQueueModule.clear();
        ad(1);

        ArrayQueueModule.dequeue();
        //ArrayQueueModule.clear();
        System.out.println(ArrayQueueModule.size());
        ad(5);
        //ArrayQueueModule.clear();
        //ArrayQueueModule.clear();
        //System.out.println(ArrayQueueModule.size());

        delAl();
        System.out.println(ArrayQueueModule.size());///size breake

        ad(5);

    }
    public static void ad(int a) {
        for (int i = 0; i < a; i++) {
            ArrayQueueModule.enqueue("e" + i);
            System.out.println(ArrayQueueModule.size());//
        }
    }
    public static void delAl() {
        while (!ArrayQueueModule.isEmpty()) {
            System.out.println(ArrayQueueModule.size() + " " + ArrayQueueModule.dequeue());
        }
    }
}
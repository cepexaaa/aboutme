package queue;

import static queue.ArrayQueueADT.*;

class ArrayQueueADTTest {
    public static void main(String[] args) {
        ArrayQueueADT q1 = create();
        ArrayQueueADT q2 = create();
        //ArrayQueueModule.clear();
        //System.out.println(ArrayQueueModule.size());
        ad(q1, 10);
        ad(q2, 5);
        delAl(q2);
        delAl(q1);
    }
    public static void ad(ArrayQueueADT q1, int a) {
        for (int i = 0; i < a; i++) {
            ArrayQueueADT.enqueue(q1, "e" + i);
            System.out.println(ArrayQueueADT.size(q1));//
        }
    }
    public static void delAl(ArrayQueueADT q1) {
        while (!ArrayQueueADT.isEmpty(q1)) {
            System.out.println(ArrayQueueADT.size(q1) + " " + ArrayQueueADT.dequeue(q1));
        }
    }
}
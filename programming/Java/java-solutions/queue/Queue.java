package queue;
// Model: a[1]..a[n]
// Inv: n >= 0 && forall i=1..n: a[i] != null
// Let: immutable(k): forall i=1..k: a'[i] = a[i]
public interface Queue {

    // Pred: element != null
    // Post: n' = n + 1 && a'[n'] = element && immutable(n)
    void enqueue(final Object x);

    //Pred: n >= 1 && queue != null
    //Post: R == a[0] && immutable(n) && n' = n
    Object element();

    //Pred: n >= 1
    //Post: n' == n - 1 && immutable(n') && R = a[0]
    Object dequeue();

    //Pred: true
    //Post: R == n && n' == n && immutable(n)
    int size();

    //Pred: true
    //Post: R == (n == 0) && n' == n && immutable(n)
    boolean isEmpty();

    //Pred: true
    //Post: n' == 0
    void clear();

    //Pred: n>=0
    //Post: n' <= n && [0, queue.length) Э i,j : i != j -> queue[i] != queue[j] &&
    //forall i = [0, queue.length) queue[i] in queue' &&
    //(any x in queue ~ i : queue[i] == x && i < j if queue[j] == x (if i != j)) == (*) &&
    //any i, j : queue[i] != queue[j] if i < j in queue => i < j in queue'
    void distinct();

}

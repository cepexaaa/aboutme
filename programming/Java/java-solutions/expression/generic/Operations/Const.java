package expression.generic.Operations;

import expression.generic.SomeExpressions;

public class Const<T> implements SomeExpressions<T> {
    private final T constt;
    public Const(final T constt) {this.constt = constt;}

    @Override
    public String toString() {
        return constt.toString();
    }

   public T evaluate(final T x) {
       return constt;
   }

    @Override
    public T evaluate(T x, T y, T z) {
        return constt;
    }

    @Override
    public int hashCode() {
        return (int) constt;
    }

}

package expression.generic;

public abstract class BinaryOperation<T> implements SomeExpressions<T> {
    private final SomeExpressions<T> argument1, argument2;
    protected final InterfaceT<T> type;
    private final String operationSign;
    public BinaryOperation(SomeExpressions<T> argument1, SomeExpressions<T> argument2, String operationSign, InterfaceT<T> type) {
        this.argument1 = argument1;
        this.argument2 = argument2;
        this.operationSign = operationSign;
        this.type = type;
    }

    public abstract T counting(T argument1, T argument2);
    public abstract int ownCode();


    public String toString() {
        return "(" + argument1 + " " + operationSign + " " + argument2 + ")";
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return counting(argument1.evaluate(x, y, z), argument2.evaluate(x, y, z));
    }


    @Override
    public int hashCode() {
        return (argument1.hashCode() * 211 + argument2.hashCode() * 179) * ownCode();
    }

}

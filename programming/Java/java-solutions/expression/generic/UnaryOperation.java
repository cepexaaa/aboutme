package expression.generic;
public abstract class UnaryOperation<T> implements SomeExpressions<T> {
    private final SomeExpressions<T> argument1;
    protected final InterfaceT<T> type;
    private final String operationSign;
    public UnaryOperation(SomeExpressions<T> argument1, String operationSign, InterfaceT<T> type) {
        this.argument1 = argument1;
        this.operationSign = operationSign;
        this.type = type;
    }

    public abstract T counting(T argument1);
    public abstract int ownCode();


    public String toString() {
        return operationSign + "(" + argument1 + ")";
    }



    public T evaluate(T x, T y, T z) {
        return counting(argument1.evaluate(x, y, z));
    }


    @Override
    public int hashCode() {
        return (argument1.hashCode() * 211) * ownCode();
    }

}

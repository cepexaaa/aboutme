package expression.generic.Operations;

import expression.generic.InterfaceT;
import expression.generic.SomeExpressions;
import expression.generic.UnaryOperation;

public class UnaryMines<T> extends UnaryOperation<T> {
    public UnaryMines(SomeExpressions<T> argument, InterfaceT<T> type) {
        super(argument, "-", type);
    }
    @Override
    public T counting(T argument) {
        return type.negate(argument);
    }
    @Override
    public int ownCode() {
        return 1131;
    }
}
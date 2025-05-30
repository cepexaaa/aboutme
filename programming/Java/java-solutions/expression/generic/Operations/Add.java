package expression.generic.Operations;

import expression.generic.BinaryOperation;
import expression.generic.InterfaceT;
import expression.generic.SomeExpressions;
public class Add<T> extends BinaryOperation<T> {
    public Add(SomeExpressions<T> argument1, SomeExpressions<T> argument2, InterfaceT<T> type) {
        super(argument1, argument2, "+", type);
    }
    @Override
    public T counting(T argument1, T argument2) {
        return type.add(argument1, argument2);
    }

    @Override
    public int ownCode() {
        return 7127;
    }

}

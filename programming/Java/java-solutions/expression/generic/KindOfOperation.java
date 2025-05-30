package expression.generic;


import expression.generic.Operations.Add;
import expression.generic.Operations.Divide;
import expression.generic.Operations.Multiply;
import expression.generic.Operations.Subtract;

public class KindOfOperation<T> {
    public SomeExpressions<T> kindOfOperation(char sign, SomeExpressions<T> arg1, SomeExpressions<T> arg2, InterfaceT<T> type) {
        return switch (sign) {
            case '+' -> new Add<>(arg1, arg2, type);
            case '-' -> new Subtract<>(arg1, arg2, type);
            case '/' -> new Divide<>(arg1, arg2, type);
            default -> new Multiply<>(arg1, arg2, type);
        };
    }
}

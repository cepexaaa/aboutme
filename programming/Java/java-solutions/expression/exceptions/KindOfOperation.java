package expression.exceptions;


public class KindOfOperation {
    public SomeExpressions kindOfOperation(char sign, SomeExpressions arg1, SomeExpressions arg2) {
        return switch (sign) {
            case '*' -> new CheckedMultiply(arg1, arg2);
            case '+' -> new CheckedAdd(arg1, arg2);
            case '-' -> new CheckedSubtract(arg1, arg2);
            case '/' -> new CheckedDivide(arg1, arg2);
//            case '&' -> new BitAnd(arg1, arg2);
//            case '|' -> new BitOr(arg1, arg2);
//            case '^' -> new BitXor(arg1, arg2);
            default -> new CheckedMultiply(arg1, arg2);
        };
    }
}

package expression.exceptions;

//import expression.SomeExpressions;

public class CheckedDivide extends BinaryOperation{
    public CheckedDivide(SomeExpressions argument1, SomeExpressions argument2) {
        super(argument1, argument2, "/");
    }
    @Override
    public int counting(int argument1, int argument2) {
        if (argument2==0) {
            throw new DivisionByZeroException("division by zero");
        } else if (argument1==Integer.MIN_VALUE && argument2 == -1) {
            throw new OverflowException("overflow");
        } else {
            return argument1 / argument2;
        }
    }

    @Override
    public int ownCode() {
        return 9013;
    }
}

package expression.exceptions;

public class CheckedAdd extends BinaryOperation {
    public CheckedAdd(SomeExpressions argument1, SomeExpressions argument2) {
        super(argument1, argument2, "+");
    }
    @Override
    public int counting(int argument1, int argument2) {
        int r = argument1 + argument2;
        if (((argument1 ^ r) & (argument2 ^ r)) >= 0) {
            return argument1 + argument2;
        }else {
            throw new OverflowException("overflow");
        }
    }

    @Override
    public int ownCode() {
        return 7127;
    }
}

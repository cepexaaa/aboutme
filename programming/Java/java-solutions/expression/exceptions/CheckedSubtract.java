package expression.exceptions;


public class CheckedSubtract extends BinaryOperation {
    public CheckedSubtract(SomeExpressions argument1, SomeExpressions argument2) {
        super(argument1, argument2, "-");
    }
    @Override
    public int counting(int argument1, int argument2) throws OverflowException {
        if ((argument2 > 0 && argument1 < Integer.MIN_VALUE + argument2) || (argument2 < 0 && argument1 > Integer.MAX_VALUE + argument2)) {
            throw new OverflowException("overflow");
        } else {
            return argument1 - argument2;
        }
    }

    @Override
    public int ownCode() {
        return 8167;
    }
}

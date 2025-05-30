package expression.exceptions;


public class CheckedNegate extends UnaryOperation {
    public CheckedNegate(SomeExpressions argument) {
        super(argument, "-");
    }
    @Override
    public int counting(int argument) throws OverflowException {
        if (argument != Integer.MIN_VALUE) {
            return argument * (-1);
        } else {
            throw new OverflowException("overflow");
        }
    }

    @Override
    public int ownCode() {
        return 1131;
    }
}

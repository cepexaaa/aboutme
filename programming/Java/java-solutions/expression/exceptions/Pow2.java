package expression.exceptions;


public class Pow2 extends UnaryOperation {
    public Pow2(SomeExpressions argument) {
        super(argument, "pow2");
    }
    @Override
    public int counting(int argument) {
        if (31<=argument || argument<0) {
            throw new OverflowException("invalid argument");
        }else {
            return 1 << argument;
        }
    }

    @Override
    public int ownCode() {
        return 1311;
    }
}

package expression.exceptions;


public class Log2 extends UnaryOperation {

    public Log2(SomeExpressions argument) {
        super(argument, "log2");
    }
    @Override
    public int counting(int argument) {
        if (argument<1) {
            throw new OverflowException("invalid argument");
        }
        int c = 0;
        while (argument > 1) {
            c++;
            argument = argument >> 1;
        }
        return c;
    }
    @Override
    public int ownCode() {
        return 3111;
    }
}

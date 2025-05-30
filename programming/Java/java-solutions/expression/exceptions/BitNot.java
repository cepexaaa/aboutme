package expression.exceptions;

public class BitNot extends UnaryOperation {
    public BitNot(SomeExpressions argument) {
        super(argument, "~");
    }
    @Override
    public int counting(int argument) {
        return ~argument;
    }
    @Override
    public int ownCode() {
        return 1111;
    }

}

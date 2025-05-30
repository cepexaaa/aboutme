package expression;

import expression.exceptions.SomeExpressions;

public class Log2 extends UnaryOperation {

    public Log2(SomeExpressions argument) {
        super(argument, "log2");
    }
    @Override
    public int counting(int argument) {
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

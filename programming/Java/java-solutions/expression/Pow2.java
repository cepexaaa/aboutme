package expression;

import expression.exceptions.SomeExpressions;

public class Pow2 extends UnaryOperation {
    public Pow2(SomeExpressions argument) {
        super(argument, "pow2");
    }
    @Override
    public int counting(int argument) {
        return 2 << argument;
    }
    @Override
    public int ownCode() {
        return 1311;
    }
}

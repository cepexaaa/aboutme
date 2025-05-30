package expression;

import expression.exceptions.SomeExpressions;

public class UnaryMines extends UnaryOperation {
    public UnaryMines(SomeExpressions argument) {
        super(argument, "-");
    }
    @Override
    public int counting(int argument) {
        return argument * (-1);
    }
    @Override
    public int ownCode() {
        return 1131;
    }
}
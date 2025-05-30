package expression.exceptions;


public class UnaryMines extends UnaryOperation {
    public UnaryMines(SomeExpressions argument) {
        super(argument, "-");
    }
    @Override
    public int counting(int argument) {
        if (argument == Integer.MIN_VALUE) {
            throw new OverflowException("overflow");
        } else {
            return argument * (-1);
        }
        
    }
    @Override
    public int ownCode() {
        return 1131;
    }
    
}
package expression.exceptions;

import java.util.List;
import java.util.Objects;

public abstract class BinaryOperation implements SomeExpressions {
    private final SomeExpressions argument1, argument2;
    private final String operationSign;
    public BinaryOperation(SomeExpressions argument1, SomeExpressions argument2, String operationSign) {
        this.argument1 = argument1;
        this.argument2 = argument2;
        this.operationSign = operationSign;
    }

    public abstract int counting(int argument1, int argument2);
    public abstract int ownCode();


    public String toString() {
        return "(" + argument1 + " " + operationSign + " " + argument2 + ")";
    }

    @Override
    public boolean equals(Object obj){
        return this == obj || (obj instanceof BinaryOperation &&
                Objects.equals(((BinaryOperation) obj).argument1, this.argument1) &&
                Objects.equals(((BinaryOperation) obj).argument2, this.argument2) &&
                Objects.equals(((BinaryOperation) obj).operationSign, this.operationSign));

    }

    @Override
    public int evaluate(int x, int y, int z)  {
        return counting(argument1.evaluate(x, y, z), argument2.evaluate(x, y, z));
    }

    @Override
    public int evaluate(List<Integer> list)  {
        return counting(argument1.evaluate(list), argument2.evaluate(list));
    }
   @Override
   public int evaluate(int x) throws DivisionByZeroException, OverflowException {
       return counting(argument1.evaluate(x), argument2.evaluate(x));
   }

    @Override
    public int hashCode() {
        return (argument1.hashCode() * 211 + argument2.hashCode() * 179) * ownCode();
    }

}

package expression;


import java.util.Objects;

import expression.exceptions.SomeExpressions;

public abstract class BinaryOperation implements SomeExpressions{
    private SomeExpressions argument1, argument2;
    private String operationSign;
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
    public int evaluate(int x, int y, int z) {
        return counting(argument1.evaluate(x, y, z), argument2.evaluate(x, y, z));
    }
   @Override
   public int evaluate(int x) {
       return counting(argument1.evaluate(x), argument2.evaluate(x));
   }

    @Override
    public int hashCode() {
        return (argument1.hashCode() * 211 + argument2.hashCode() * 179) * ownCode();
    }

}

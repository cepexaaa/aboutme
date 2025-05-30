package expression.exceptions;

import java.util.List;
import java.util.Objects;

public abstract class UnaryOperation implements SomeExpressions {
    private final SomeExpressions argument1;
    private final String operationSign;
    public UnaryOperation(SomeExpressions argument1, String operationSign) {
        this.argument1 = argument1;
        this.operationSign = operationSign;
    }

    public abstract int counting(int argument1);
    public abstract int ownCode();


    public String toString() {
        return operationSign + "(" + argument1 + ")";
    }

    @Override
    public boolean equals(Object obj){
        return this == obj || (obj instanceof UnaryOperation &&
                Objects.equals(((UnaryOperation) obj).argument1, this.argument1) &&
                Objects.equals(((UnaryOperation) obj).operationSign, this.operationSign));

    }

    @Override
    public int evaluate(int x, int y, int z) {
        return counting(argument1.evaluate(x, y, z));
    }
   @Override
   public int evaluate(int x) {
       return counting(argument1.evaluate(x));
   }

    @Override
    public int hashCode() {
        return (argument1.hashCode() * 211) * ownCode();
    }
    @Override
    public int evaluate(List<Integer> variables) {
        return counting(argument1.evaluate(variables));
//        argument1.evaluate(variables);
//        throw new UnsupportedOperationException("Unimplemented method 'evaluate'");
    }

}

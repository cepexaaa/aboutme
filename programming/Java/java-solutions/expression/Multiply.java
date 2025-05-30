package expression;

import java.util.List;

import expression.exceptions.SomeExpressions;

public class Multiply extends BinaryOperation {
    public Multiply(SomeExpressions argument1, SomeExpressions argument2) {
        super(argument1, argument2, "*");
    }
    @Override
    public int counting(int argument1, int argument2) {
        return argument1 * argument2;
    }

    @Override
    public int ownCode() {
        return 5171;
    }
       @Override
   public int evaluate(List<Integer> list) {
        return 0;
   }
}

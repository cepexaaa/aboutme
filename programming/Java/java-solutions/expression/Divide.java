package expression;

import java.util.List;

import expression.exceptions.SomeExpressions;

public class Divide extends BinaryOperation{
    public Divide(SomeExpressions argument1, SomeExpressions argument2) {
        super(argument1, argument2, "/");
    }
    @Override
    public int counting(int argument1, int argument2) {
        return argument1 / argument2;
    }

    @Override
    public int ownCode() {
        return 9013;
    }
       @Override
   public int evaluate(List<Integer> list) {
        return 0;
   }
   
}

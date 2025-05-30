package expression;

import java.util.List;

import expression.exceptions.SomeExpressions;

public class Add extends BinaryOperation {
    public Add(SomeExpressions argument1, SomeExpressions argument2) {
        super(argument1, argument2, "+");
    }
    @Override
    public int counting(int argument1, int argument2) {
        return argument1 + argument2;
    }

    @Override
    public int ownCode() {
        return 7127;
    }
       @Override
   public int evaluate(List<Integer> list) {
        return 0;
   }
}

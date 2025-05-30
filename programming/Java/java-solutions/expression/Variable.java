package expression;

import java.util.List;
import java.util.Objects;

import expression.exceptions.SomeExpressions;

public class Variable implements SomeExpressions {
    private String value;
    private int position;

    public Variable(String var) {
        this.value = var;
    }
    
    public Variable (int var) {
        this.value = null;
        this.position = var;
    }

    @Override
    public String toString() {
        if (value == null) {
            return "$" + position;
        } else {
            return value;
        }
    }

   @Override
   public int evaluate(final int x) {
           return(x);
   }

   @Override
   public int evaluate(List<Integer> list) {
        return 0;
   }

    @Override
    public int evaluate(int x, int y, int z) {
        return switch (this.value) {
            case "y" -> y;
            case "z" -> z;
            default -> x;
        };
    }

    @Override
    public int hashCode() {
        return value.hashCode();
    }

     @Override
    public boolean equals(Object obj) {
        if (obj instanceof Variable) {
            Variable o = (Variable) obj;
            return Objects.equals(value, o.value);
        } return false;
    }
}

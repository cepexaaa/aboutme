package expression.generic.Operations;

import expression.generic.SomeExpressions;

import java.util.Objects;

public class Variable<T> implements SomeExpressions<T> {
    private final String value;

    public Variable(String var) {
        this.value = var;
    }


    @Override
    public String toString() {
        return value;
    }


   public T evaluate(final T x) {
           return(x);
   }

    @Override
    public T evaluate(T x, T y, T z) {
        assert value != null;
        return switch (value) {
            case "y" -> y;
            case "z" -> z;
            default -> x;
        };
    }

    @Override
    public int hashCode() {
        return Objects.hash(value);
    }
}

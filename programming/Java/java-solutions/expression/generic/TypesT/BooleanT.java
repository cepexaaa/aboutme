package expression.generic.TypesT;

import expression.exceptions.DivisionByZeroException;
import expression.generic.InterfaceT;

public class BooleanT implements InterfaceT<Boolean> {

    @Override
    public Boolean add(Boolean first, Boolean second) {
        return first || second;
    }

    @Override
    public Boolean substract(Boolean first, Boolean second) {
        return first ^ second;
    }

    @Override
    public Boolean divide(Boolean first, Boolean second) {
        if (!second) {
            throw new DivisionByZeroException("division by zero");
        }
        return first;
    }

    @Override
    public Boolean multipy(Boolean first, Boolean second) {
        return first && second;
    }

    @Override
    public Boolean negate(Boolean first) {
        return first;
    }

    @Override
    public Boolean valueOf(int x) {
        if (x != 0) {
            return true;
        }
        return Boolean.valueOf(String.valueOf(x));
    }

    @Override
    public Boolean parse(String x) {
        return !x.equals("0");
    }
}

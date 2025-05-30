package expression.generic.TypesT;

import expression.exceptions.OverflowException;
import expression.generic.InterfaceT;

public class IntegerT implements InterfaceT<Integer> {
    private final boolean isOverflowCheck;

    public IntegerT(boolean isOverflowCheck) {
        this.isOverflowCheck = isOverflowCheck;
    }

    @Override
    public Integer add(Integer first, Integer second) {
        if (isOverflowCheck) {
            if ((first > 0 && second > 0 && first > Integer.MAX_VALUE - second) || (first < 0 && second < 0 && first < Integer.MIN_VALUE - second)) {
                throw new OverflowException("overflow with +");
            } else {
                return first + second;
            }
        }
        return first + second;
    }

    @Override
    public Integer substract(Integer first, Integer second) {
        if (isOverflowCheck) {
            if ((second > 0 && first < Integer.MIN_VALUE + second) || (second < 0 && first > Integer.MAX_VALUE + second)) {
                throw new OverflowException("overflow with '-'");
            } else {
                return first - second;
            }
        }
        return first - second;
    }

    @Override
    public Integer divide(Integer first, Integer second) {
        if (isOverflowCheck) {
            if (first==Integer.MIN_VALUE && second == -1) {
                throw new OverflowException("overflow");
            }
            return first / second;
        }
        return first / second;
    }

    @Override
    public Integer multipy(Integer first, Integer second) {
        if (isOverflowCheck) {
            if (isMultiplicationOverflow(first, second)) {
                throw new OverflowException("overflow");
            } return  first * second;
        }
        return first * second;
    }

    @Override
    public Integer negate(Integer first) {
        if (isOverflowCheck) {
            if (first != Integer.MIN_VALUE) {
                return first * (-1);
            } else {
                throw new OverflowException("overflow");
            }
        }
        return -first;
    }

    @Override
    public Integer valueOf(int x) {
        return x;
    }
    public static boolean isMultiplicationOverflow(int a, int b) {
        if (a > 0) {
            if (b > 0) {
                return a > Integer.MAX_VALUE / b;
            } else if (b < 0) {
                return b < Integer.MIN_VALUE / a;
            }
        } else if (a < 0) {
            if (b > 0) {
                return a < Integer.MIN_VALUE / b;
            } else if (b < 0) {
                return (b < Integer.MAX_VALUE / a);
            }
        }
        return false;
    }
    @Override
    public Integer parse(String s) {
        return Integer.parseInt(s);
    }
}
package expression.exceptions;

//import expression.TripleExpression;

public class CheckedMultiply extends BinaryOperation {
    public CheckedMultiply(SomeExpressions argument1, SomeExpressions argument2) {
        super(argument1, argument2, "*");
    }
    @Override
    public int counting(int argument1, int argument2) {
        //if ((argument2 != 0) && (argument1 > Integer.MAX_VALUE / argument2) || (argument1 == Integer.MIN_VALUE) && (argument2 == -1)){
        if (isMultiplicationOverflow(argument1, argument2)) {
            throw new OverflowException("overflow");
        } else {
            //System.err.println(argument1 + " " + argument2);
            return argument1 * argument2;
        }
    }

    @Override
    public int ownCode() {
        return 5171;
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
                return (a != 0) && (b < Integer.MAX_VALUE / a);
            }
        // }if (a == -2147483648 && b == -2147483648) {
        //     return true;
        }
        return false;

        //return true;
        // long rez = a * b;
        // if ((int) rez != a * b) {
        //     return true;
        // } else {
        //     return false;
        // }
    }
}

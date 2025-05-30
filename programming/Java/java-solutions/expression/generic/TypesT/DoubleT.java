package expression.generic.TypesT;

import expression.generic.InterfaceT;

public class DoubleT implements InterfaceT<Double> {

    @Override
    public Double add(Double first, Double second) {
        return first + second;
    }

    @Override
    public Double substract(Double first, Double second) {
        return first - second;
    }

    @Override
    public Double divide(Double first, Double second) {
        return first / second;
    }

    @Override
    public Double multipy(Double first, Double second) {
        return first * second;
    }

    @Override
    public Double negate(Double first) {
        return -first;
    }

    @Override
    public Double valueOf(int x) {
        return (double) x;
    }

    @Override
    public Double parse(String x) {
        return Double.parseDouble(x);
    }

}
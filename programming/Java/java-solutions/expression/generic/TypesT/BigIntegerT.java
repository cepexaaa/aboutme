package expression.generic.TypesT;

import expression.generic.InterfaceT;

import java.math.BigInteger;

public class BigIntegerT implements InterfaceT<BigInteger> {
    @Override
    public BigInteger add(BigInteger first, BigInteger second) {
        return first.add(second);
    }

    @Override
    public BigInteger substract(BigInteger first, BigInteger second) {
        return first.subtract(second);
    }

    @Override
    public BigInteger divide(BigInteger first, BigInteger second) {
        return first.divide(second);
    }

    @Override
    public BigInteger multipy(BigInteger first, BigInteger second) {
        return first.multiply(second);
    }

    @Override
    public BigInteger negate(BigInteger first) {
        return first.negate();
    }

    @Override
    public BigInteger valueOf(int x) {
        return BigInteger.valueOf(x);
    }

    @Override
    public BigInteger parse(String x) {
        return new BigInteger(x);
    }
}
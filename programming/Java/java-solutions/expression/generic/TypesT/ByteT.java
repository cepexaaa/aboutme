package expression.generic.TypesT;

import expression.generic.InterfaceT;

public class ByteT implements InterfaceT<Byte> {
    @Override
    public Byte add(Byte first, Byte second) {
        return (byte) (first + second);
    }

    @Override
    public Byte substract(Byte first, Byte second) {
        return (byte) (first - second);
    }

    @Override
    public Byte divide(Byte first, Byte second) {
        return (byte) (first / second);
    }

    @Override
    public Byte multipy(Byte first, Byte second) {
        return (byte) (first * second);
    }

    @Override
    public Byte negate(Byte first) {
        return (byte)(-first);
    }

    @Override
    public Byte valueOf(int x) {
        return (byte) x;
    }

    @Override
    public Byte parse(String x) {
        return Byte.parseByte(x);
    }

}

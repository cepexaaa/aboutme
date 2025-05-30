package expression.generic;

public interface InterfaceT<T> {
    T add(T first, T second);
    T substract(T first, T second);
    T divide(T first, T second);
    T multipy(T first, T second);
    T negate(T first);
    T valueOf(int x);
    T parse(String x);
}

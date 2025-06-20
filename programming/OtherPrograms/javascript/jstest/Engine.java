package jstest;

import base.Asserts;

import java.util.function.BiFunction;

/**
 * Test engine.
 *
 * @author Georgiy Korneev (kgeorgiy@kgeorgiy.info)
 */
public interface Engine<X> {
    Result<X> prepare(String expression);

    Result<Number> evaluate(final Result<X> prepared, double[] vars);

    Result<String> toString(final Result<X> prepared);

    default Result<X> parse(final String expression) {
        throw new UnsupportedOperationException();
    }

    record Result<T>(String context, T value) {
        public void assertEquals(final T expected) {
            Asserts.assertEquals(context(), expected, value());
        }

        public <R> Result<R> cast(final BiFunction<T, String, R> convert) {
            return new Result<>(context(), convert.apply(value(), context()));
        }

        @Override
        public String toString() {
            return context();
        }
    }
}

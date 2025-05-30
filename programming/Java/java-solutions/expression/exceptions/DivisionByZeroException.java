package expression.exceptions;

public class DivisionByZeroException extends RuntimeException{
    public DivisionByZeroException(String msg) {
        super(msg);
    }

    public DivisionByZeroException(String message, Throwable cause) {
        super(message, cause);
    }
}

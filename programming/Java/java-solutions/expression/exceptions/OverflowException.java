package expression.exceptions;

public class OverflowException extends RuntimeException{

    public OverflowException(String msg) {
        super(msg);
    }

    public OverflowException(String message, Throwable cause) {
        super(message, cause);
    }
}

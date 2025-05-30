package expression.exceptions;

public class NoArgumentException extends ExcetpionsOfParser{
    public NoArgumentException(String a){
        super(a);
    }
    public NoArgumentException(String message, Throwable cause) {
        super(message, cause);
    }
}

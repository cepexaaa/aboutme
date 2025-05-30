package expression.exceptions;

public class ParenthesisCountException extends ExcetpionsOfParser{
    public ParenthesisCountException(String a){
        super(a);
    }
    public ParenthesisCountException(String message, Throwable cause) {
        super(message, cause);
    }
}

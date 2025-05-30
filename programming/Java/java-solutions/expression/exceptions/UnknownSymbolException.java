package expression.exceptions;

public class UnknownSymbolException extends ExcetpionsOfParser{
    public UnknownSymbolException(String a){
        super(a);
    }
    public UnknownSymbolException(String message, Throwable cause) {
        super(message, cause);
    }
}

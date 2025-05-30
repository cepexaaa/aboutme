package expression.exceptions;

import java.util.List;


public class Const implements SomeExpressions {
    private int constt;
    public Const(final int constt) {this.constt = constt;}

    @Override
    public String toString() {
        return Integer.toString(constt);
    }

    @Override
   public int evaluate(final int x) {
       return constt;
   }

    @Override
    public int evaluate(int x, int y, int z) {
        return constt;
    }

    @Override
    public int hashCode() {
        return constt;
    }

    @Override
    public boolean equals(Object object) {
        if (object instanceof Const) {
            Const compared = (Const) object;
            return this.constt == compared.constt;
        }
        return false;
    }
     @Override
    public int evaluate(List<Integer> variables) {
        return constt;
    }


}

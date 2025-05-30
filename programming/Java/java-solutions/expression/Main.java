package expression;

import java.util.List;

import expression.exceptions.ExcetpionsOfParser;
import expression.exceptions.ExpressionParser;

public class Main {
    public static void main(String[] args) throws ExcetpionsOfParser {
        for (int i = 0; i < 10; i++) {
            System.out.println(new ExpressionParser().parse("1000000*x*x*x*x*x/(x-1)"));    
        }
    }
}

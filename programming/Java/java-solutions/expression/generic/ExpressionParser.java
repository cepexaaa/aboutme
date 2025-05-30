package expression.generic;

import expression.generic.Operations.Const;
import expression.generic.Operations.UnaryMines;
import expression.generic.Operations.Variable;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ExpressionParser<T>{
    private final static Map<Character, Byte> PRIORITIES = new HashMap<>();
    public ExpressionParser(InterfaceT<T> type) {
        this.type = type;
    }
    private final InterfaceT<T> type;
    private int input;
    private char c;
    private boolean flagCounting;
    private String expression;

    public SomeExpressions<T> parse(String expression){
        PRIORITIES.put('*', (byte) 1);
        PRIORITIES.put('/', (byte) 1);
        PRIORITIES.put('-', (byte) 2);
        PRIORITIES.put('+', (byte) 2);
        input = 0;
        flagCounting = true;
        this.expression = expression.trim();
        return parsing();
       
    }

    public SomeExpressions<T> parsing(){
        boolean flagUnary = true;
        int end = expression.length();
        List<Byte> stackOfOperationPriority = new ArrayList<>();
        List<Character> stackOfOperationSign = new ArrayList<>();
        List<SomeExpressions<T>> stackOfValues = new ArrayList<>();
        KindOfOperation<T> operation = new KindOfOperation<T>();
        boolean flag;
        while (input < end) {
            c = expression.charAt(input++);
            
            if (Character.isWhitespace(c) || c == '\r' || c == '\n') {
                continue;
            }
            else if (c == 'x' || c == 'y' || c == 'z') {
                stackOfValues.add(new Variable<>(Character.toString(c)));
                flagUnary = false;
                continue;
            }
            else if (flagUnary && c == '-') {
                input--;
                stackOfValues.add(collectingAnUnaryExpression());
                flagUnary = false;
                if (input >= end && c != ')') {
                    break;
                }
            }
            String  number = null;
            flag = false;
            
            if (Character.isDigit(c) && flagCounting) {
                flag = true;
                number = wholeNumber(false);
            }
            if (flag) {
                stackOfValues.add(new Const<>(type.parse(number)));
                flagUnary = false;
            }
            if (PRIORITIES.containsKey(c)) {
                while (!stackOfOperationSign.isEmpty()) {
                    byte priority = stackOfOperationPriority.get(stackOfOperationSign.size() - 1);
                    if (priority <= PRIORITIES.get(c)) {
                        SomeExpressions<T> operation2 = stackOfValues.remove(stackOfValues.size() - 1);
                        SomeExpressions<T> operation1 = stackOfValues.remove(stackOfValues.size() - 1);
                        char sign = stackOfOperationSign.remove(stackOfOperationSign.size() - 1);
                        stackOfOperationPriority.remove(stackOfOperationPriority.size() - 1);
                        stackOfValues.add(operation.kindOfOperation(sign, operation1, operation2, type));
                    } else {
                        break;
                    }
                }
                stackOfOperationPriority.add(PRIORITIES.get(c));
                stackOfOperationSign.add(c);
                flagUnary = true;
            }
            else if (c == ')') {
                return calc(stackOfValues, stackOfOperationSign, operation);
            }

            else if (c == '(') {
                stackOfValues.add(parsing());
                flagUnary = false;
            }
        }
        return calc(stackOfValues, stackOfOperationSign, operation);
    }

    private SomeExpressions<T> calc(List<SomeExpressions<T>> stackOfValues, List<Character> stackOfOperationSign, KindOfOperation<T> operation) {
        while (stackOfValues.size() > 1) {
            SomeExpressions<T> operation2 = stackOfValues.remove(stackOfValues.size() - 1);
            SomeExpressions<T> operation1 = stackOfValues.remove(stackOfValues.size() - 1);
            char sn = stackOfOperationSign.remove(stackOfOperationSign.size()-1);
            stackOfValues.add(operation.kindOfOperation(sn, operation1, operation2, type));
        }
        return stackOfValues.remove(0);
    }

    public SomeExpressions<T> collectingAnUnaryExpression(){
        
        int end = expression.length();
        while (input < end) {
            c = expression.charAt(input++);
                if (c == ' ') {
                    continue;
                } else if (c == 'x' || c == 'y' || c == 'z') {
                    return new Variable<>(Character.toString(c));
                }
                else if (c == '-'){
                    while (expression.charAt(input) == ' ') {
                        input++;
                    }
                    if (Character.isDigit(expression.charAt(input))) {
                        c = expression.charAt(input++);
                        return new Const<>(type.parse(wholeNumber(true)));
                    } else {
                        return new UnaryMines<>(collectingAnUnaryExpression(), type);
                    }
                }
                else if (c == '('){
                    SomeExpressions<T> e = parsing();
                    if (input < end) {
                        c = expression.charAt(input++);
                    }
                    return e;
                } else {
                    if (Character.isDigit(c)) {
                        return new Const<>(type.parse(wholeNumber(true)));
                    }
                }
            }
        throw new RuntimeException();
    }

    public String wholeNumber(boolean m) {
        int end = expression.length();
        StringBuilder number = new StringBuilder();
        if (m) {
            number.append("-");
        }
        while (Character.isDigit(c) && flagCounting) {
            number.append(c);
            if (input < end) {
                c = expression.charAt(input++);
            } else {
                flagCounting = false;
                break;
            }
        }
        return number.toString();
        //return Integer.parseInt(number.toString());
    }

}

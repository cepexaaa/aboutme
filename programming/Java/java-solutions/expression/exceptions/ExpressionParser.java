package expression.exceptions;

import expression.*;

import java.util.*;

public class ExpressionParser implements expression.exceptions.TripleParser, ListParser {
    static Map<Character, Byte> PRIORITIES = new HashMap<>() {{
        put('*', (byte) 1);
        put('+', (byte) 2);
        put('-', (byte) 2);
        put('/', (byte) 1);
        put('^', (byte) 4);
        put('&', (byte) 3);
        put('|', (byte) 5);
    }};
    private int input;
    private char c;
    private boolean flagCounting;
    private String expression;
    private int parenthesisCounter;
    private List<Character> balance;

    public ListExpression parse(String expression, List<String> listVar) throws ExcetpionsOfParser {
        input = 0;
        flagCounting = true;
        balance = new ArrayList<>();
        this.expression = expression.trim();
        parenthesisCounter = 0;
        SomeExpressions rezult = parsing(listVar);
        if (parenthesisCounter<0) {
            throw new ParenthesisCountException("No opening parenthesis " + expression);
        } else if (parenthesisCounter>0) {
            throw  new ParenthesisCountException("No closing parenthesis " + expression);
        } else return rezult;
    }

    public SomeExpressions parse(String expression) throws ExcetpionsOfParser {
        input = 0;
        flagCounting = true;
        balance = new ArrayList<>();
        this.expression = expression.trim();
        parenthesisCounter = 0;
        SomeExpressions rezult1 = parsing(null);
        if (parenthesisCounter<0) {
            throw new ParenthesisCountException("No opening parenthesis " + input + " " + expression + " " + parenthesisCounter);
        } else if (parenthesisCounter>0) {
            throw  new ParenthesisCountException("No closing parenthesis " + input + " " + expression + " " + parenthesisCounter);
        } else return rezult1;
       
    }

    public SomeExpressions parsing(List<String> listVar) throws ExcetpionsOfParser {
        boolean flagUnary = true;
        int end = expression.length();
        List<Byte> stackOfOperationPriority = new ArrayList<>();
        List<Character> stackOfOperationSign = new ArrayList<>();
        List<SomeExpressions> stackOfValues = new ArrayList<>();
        KindOfOperation operation = new KindOfOperation();
        boolean flag;
        int ans;
        
        while (input < end) {
            c = expression.charAt(input++);
            
            if (Character.isWhitespace(c) || c == '\r' || c == '\n') {
                continue;
            }
            else if (c == 'x' || c == 'y' || c == 'z') {
                stackOfValues.add(new Variable(Character.toString(c)));
                flagUnary = false;
                continue;
            } else if ((ans = isSomeWord(c, listVar)) > -1){
                stackOfValues.add(new Variable(ans));
                flagUnary = false;
                
                if (input >= end) {
                    break;
                } else {
                    c = expression.charAt(input++);
                }
            }
            else if (flagUnary && (c == '~' || c == '-' || c == 'l' || c == 'p')) {
                input--;
                stackOfValues.add(collectingAnUnaryExpression(listVar));
                flagUnary = false;
                if (input >= end && c != ')' && c != ']' && c != '}') {
                    break;
                }
            }
            int number = 0;
            flag = false;
            
            if (Character.isDigit(c) && flagCounting) {
                
                if (!flagUnary) {
                    throw new SpaceInNumberException("Space in number" + c + expression);
                }
                flag = true;
                number = wholeNumber(false);
            }
            if (flag) {
                stackOfValues.add(new Const(number));
                flagUnary = false;
            }
            if (PRIORITIES.containsKey(c)) {
                if (flagUnary && c != '-') {
                    if (stackOfValues.isEmpty()){
                        throw new NoArgumentException("No first argument " + expression);
                    } else {
                        throw new NoArgumentException("No middle argument " + expression);
                    }
                    
                }
                if (input >= (end)) {
                    throw new NoArgumentException("No last argument " + expression);
                } 
                while (!stackOfOperationSign.isEmpty()) {
                    byte priority = stackOfOperationPriority.get(stackOfOperationSign.size() - 1);
                    if (priority <= PRIORITIES.get(c)) {
                        SomeExpressions operation2 = stackOfValues.remove(stackOfValues.size() - 1);
                        SomeExpressions operation1 = stackOfValues.remove(stackOfValues.size() - 1);
                        char sign = stackOfOperationSign.remove(stackOfOperationSign.size() - 1);
                        stackOfOperationPriority.remove(stackOfOperationPriority.size() - 1);
                        stackOfValues.add(operation.kindOfOperation(sign, operation1, operation2));
                    } else {
                        break;
                    }
                }
                stackOfOperationPriority.add(PRIORITIES.get(c));
                stackOfOperationSign.add(c);
                flagUnary = true;
            }

            
            else if (isCloseParenthesis(c)) {
                
                parenthesisCounter--;
                if (flagUnary){
                    throw new NoArgumentException("No last argument " + c + " " + expression);
                }
                while (stackOfValues.size() > 1) {
                    SomeExpressions operation2 = stackOfValues.remove(stackOfValues.size() - 1);
                    SomeExpressions operation1 = stackOfValues.remove(stackOfValues.size() - 1);
                    char sn = stackOfOperationSign.remove(stackOfOperationSign.size() - 1);
                    stackOfValues.add(operation.kindOfOperation(sn, operation1, operation2));
                }
                return stackOfValues.remove(0);
            }

            else if (isOpenParenthesis(c)) {
                parenthesisCounter++;
                stackOfValues.add(parsing(listVar));
                flagUnary = false;
            }
            else if ((!PRIORITIES.containsKey(c)) && (!Character.isDigit(c)) && (!Character.isWhitespace(c)) && (c != 'x') && (c != 'y') && (c != 'z')){
                boolean falgexception = true;
                int ind = 9;
                String maybeVar = null;
                if (listVar!=null) {
                    
                    maybeVar = collectIsMaybeVar();
                    ind = listVar.indexOf(maybeVar);
                
                    if (ind != -1) {
                        
                        flagUnary = false;
                        falgexception = false;
                        stackOfValues.add(new Variable(ind));
                    } 
                    
                } if (falgexception) { 
                    
                    if (input==1){
                        throw new UnknownSymbolException("Start symbol '" + c + "' in " + expression);
                    } else if (input == end) {
                        throw new UnknownSymbolException("End symbol '" + c + "' in " + expression);
                    }else {
                        throw new UnknownSymbolException("Middle symbol '" + c + "' in " + expression);
                    }
                }
                
            }
        }
        while (stackOfValues.size() > 1) {
            
            SomeExpressions operation2 = stackOfValues.remove(stackOfValues.size() - 1);
            SomeExpressions operation1 = stackOfValues.remove(stackOfValues.size() - 1);
            char sn = stackOfOperationSign.remove(stackOfOperationSign.size()-1);
            stackOfValues.add(operation.kindOfOperation(sn, operation1, operation2));
        } 
        return stackOfValues.remove(0);
    }

    public int isSomeWord(char c, List<String> listVar) throws ExcetpionsOfParser { 
        int ind = -1;
        if((!PRIORITIES.containsKey(c)) && (!Character.isDigit(c)) && (!Character.isWhitespace(c)) && c != ')' && c != '('){
            if (listVar!=null) {
                
                String maybeVar = collectIsMaybeVar();
                ind = listVar.indexOf(maybeVar);
                
                return ind;
            }
        }
        return ind;
    }

    public String collectIsMaybeVar() {
        StringBuilder sb = new StringBuilder();
        sb.append(c);
        
        while (input < expression.length() && expression.charAt(input) != ')' && !Character.isWhitespace(expression.charAt(input)) && !PRIORITIES.containsKey(expression.charAt(input))) {
            sb.append(expression.charAt(input++));
        
        }         
        return sb.toString();
    }

    public SomeExpressions collectingAnUnaryExpression(List<String> listVar) throws ExcetpionsOfParser {
        
        int end = expression.length();
        while (input < end) {
            
            c = expression.charAt(input++);
            
                if (c == ' ') {
                    continue;
                } else if (c == 'x' || c == 'y' || c == 'z') {
                    return new Variable(Character.toString(c)); 
                }
                else if (isLog2(c)) {
                    
                    while (expression.charAt(input) == ' ') {
                        input++;
                    }
                    if (expression.charAt(input) != '(' && expression.charAt(input-1) == '2') {
                        throw new OverflowException("uncorrect expression" + " " + expression);
                    }
                    if (Character.isDigit(expression.charAt(input))) {
                        c = expression.charAt(input++);
                        return new Const(wholeNumber(false));
                    } else {
                        return new Log2(collectingAnUnaryExpression(listVar));
                    }
                }
                else if (isPow2(c)) {                    
                    
                    while (expression.charAt(input) == ' ') {
                        input++;                        
                    }
                    if (expression.charAt(input) != '(' && expression.charAt(input-1) == '2') {
                        throw new OverflowException("uncorrect expression" + " " + expression);
                    }
                    if (Character.isDigit(expression.charAt(input))) {
                        
                        c = expression.charAt(input++);
                        return new Const(wholeNumber(false));
                    } else {
                        
                        return new Pow2(collectingAnUnaryExpression(listVar));
                    }
                }
                else if (c == '~'){
                    while (expression.charAt(input) == ' ') {
                        input++;
                    }
                    if (Character.isDigit(expression.charAt(input))) {
                        c = expression.charAt(input++);
                        return new Const(~wholeNumber(false));
                    } else {
                        return new BitNot(collectingAnUnaryExpression(listVar));
                    }
                }
                else if (c == '-'){
                    if (input==end) {
                        throw new NoArgumentException("No last argument " + expression);
                    }
                    while (expression.charAt(input) == ' ') {
                        input++;
                    }
                    if (Character.isDigit(expression.charAt(input))) {
                        c = expression.charAt(input++);
                        return new Const(wholeNumber(true));
                    } else {
                        return new UnaryMines(collectingAnUnaryExpression(listVar));
                    }
                }
                else if (isOpenParenthesis(c)){
                    
                    parenthesisCounter++;
                    SomeExpressions e = parsing(listVar);
                    
                    if (input < end) {
                        c = expression.charAt(input++);
                    } else {
                        parenthesisCounter++;
                    }
                    return e;
                } else {
                
                    if (Character.isDigit(c)) {
                        return new Const(wholeNumber(true));
                    } else {
                        boolean falgexception = true;
                        if (listVar!=null) {
                           
                            String maybeVar = collectIsMaybeVar();
                            int ind = listVar.indexOf(maybeVar);
                            
                            if (ind != -1) {
                                
                                
                                falgexception = false;
                                if (input<end) {
                                    c = expression.charAt(input++);
                                } else {
                                    c = expression.charAt(input-1);
                                }
                                return new Variable(ind);
                            } 
                            
                        }
                        if (PRIORITIES.containsKey(c) && falgexception) {
                            throw new NoArgumentException("No middle symbol " + expression);
                        }
                    } 
                }
            }

        
        throw new UnknownSymbolException("UnknownSymbol " + expression.charAt(input-1));
    }

    public int wholeNumber(boolean m) {
        int end = expression.length();
        StringBuilder number = new StringBuilder();
        if (m) {
            number.append("-");
        }
        while (Character.isDigit(c) && flagCounting) {
            number.append(Character.toString(c));
            if (input < end) {
                c = expression.charAt(input++);
            } else {
                flagCounting = false;
                break;
            }
        }
        return Integer.parseInt(number.toString());
    }



    public boolean isLog2(char c) {
        if (c == 'l' && expression.charAt(input) == 'o' && expression.charAt(input+1) == 'g' && expression.charAt(input+2) == '2') {
            input+=3;
            return true;
        } else{
            return false;
        } 
    }

    public boolean isPow2(char c) {
        if (c == 'p' && expression.charAt(input) == 'o' && expression.charAt(input+1) == 'w' && expression.charAt(input+2) == '2') {
            input+=3;
            return true;
        } else{
            return false;
        } 
    }



    public boolean isOpenParenthesis(char c) {
        
        switch (c) {
            case '(':
                balance.add('(');
                return true;
            case '[':
                balance.add('[');
                
                return true;
            case '{':
                balance.add('{');
                return true;        
            default:
                return false;
        }
    }

    public boolean isCloseParenthesis(char c) throws ParenthesisCountException {
        
        switch (c) {
            case ')':
            if (input >= expression.length() && balance.isEmpty()) {
                return true;
            }
                if (!balance.isEmpty() && balance.get(balance.size()-1) == '(') {
                    balance.remove(balance.size()-1);
                    return true;
                } else {
                    throw new ParenthesisCountException("No twise parenthesis " + expression);
                }
            case ']':
            if (input >= expression.length() && balance.isEmpty()) {
                return true;
            }
                
                if (!balance.isEmpty() && balance.get(balance.size()-1) == '[') {
                    balance.remove(balance.size()-1);
                    return true;
                }
                else {
                    throw new ParenthesisCountException("No twise parenthesis " + expression);
                }
            case '}':
            if (input >= expression.length() && balance.isEmpty()) {
                return true;
            }
                if (!balance.isEmpty() && balance.get(balance.size()-1) == '{') {
                    balance.remove(balance.size()-1);
                    return true;
                }
                else {
                    throw new ParenthesisCountException("No twise parenthesis " + expression);
                }
            default:
                
                return false;
        }
    }
}

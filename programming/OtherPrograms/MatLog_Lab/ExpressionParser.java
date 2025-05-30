
import java.util.*;

public class ExpressionParser {

    // Типы токенов
    enum TokenType {
        VAR, IMPL, AND, OR, NOT, LPAR, RPAR, TURN, COMMA
    }

    // Класс для представления токена
    static class Token {
        TokenType type;
        String value;

        Token(TokenType type, String value) {
            this.type = type;
            this.value = value;
        }
    }

    // Лексический анализатор
    static List<Token> tokenize(String input) {
        List<Token> tokens = new ArrayList<>();
        int i = 0;
        while (i < input.length()) {
            char c = input.charAt(i);
            if (Character.isWhitespace(c) || c == '\t' || c == '\r') {
                i++;
            } else if (c == '|') {
                tokens.add(new Token(TokenType.OR, "|"));
                i++;
            } else if (c == '&') {
                tokens.add(new Token(TokenType.AND, "&"));
                i++;
            } else if (c == '-' && i + 1 < input.length() && input.charAt(i + 1) == '>') {
                tokens.add(new Token(TokenType.IMPL, "->"));
                i += 2;
            } else if (c == '!') {
                tokens.add(new Token(TokenType.NOT, "!"));
                i++;
            } else if (c == '(') {
                tokens.add(new Token(TokenType.LPAR, "("));
                i++;
            } else if (c == ')') {
                tokens.add(new Token(TokenType.RPAR, ")"));
                i++;
            } else if (c == ',') {
                tokens.add(new Token(TokenType.COMMA, ","));
                i++;
            } else if (Character.isUpperCase(c)|| c == 39) {
                StringBuilder var = new StringBuilder();
                while (i < input.length() && (Character.isLetterOrDigit(input.charAt(i)) || input.charAt(i) == 39)) {
                    var.append(input.charAt(i));
                    i++;
                }
                tokens.add(new Token(TokenType.VAR, var.toString()));
            } else {
                i++; // Пропускаем недопустимые символы
            }
        }
        return tokens;
    }

    // Парсер для выражений
    static class Parser {
        private final List<Token> tokens;
        private int pos;

        Parser(List<Token> tokens) {
            this.tokens = tokens;
            this.pos = 0;
        }

        // Основной метод парсинга
        String parse() {
            Expr expr = parseImpl();
            return expr.toPrefixString();
        }

        // Парсинг импликации
        private Expr parseImpl() {
            Expr left = parseOr();
            if (match(TokenType.IMPL)) {
                Expr right = parseImpl();
                return new ImplExpr(left, right);
            }
            return left;
        }

        // // Парсинг дизъюнкции
        // private Expr parseOr() {
        //     Expr left = parseAnd();
        //     if (match(TokenType.OR)) {
        //         Expr right = parseOr();
        //         return new OrExpr(left, right);
        //     }
        //     return left;
        // }

        private Expr parseOr() {
            Expr left = parseAnd();
            while (match(TokenType.OR)) {
                Expr right = parseAnd();
                left = new OrExpr(left, right);
            }
            return left;
        }

        // Парсинг конъюнкции (левоассоциативная)
        private Expr parseAnd() {
            Expr left = parseNot();
            while (match(TokenType.AND)) {
                Expr right = parseNot();
                left = new AndExpr(left, right);
            }
            return left;
        }

        // // Парсинг конъюнкции
        // private Expr parseAnd() {
        //     Expr left = parseNot();
        //     if (match(TokenType.AND)) {
        //         Expr right = parseAnd();
        //         return new AndExpr(left, right);
        //     }
        //     return left;
        // }

        // Парсинг отрицания
        private Expr parseNot() {
            if (match(TokenType.NOT)) {
                Expr expr = parseNot();
                return new NotExpr(expr);
            }
            return parsePrimary();
        }

        // Парсинг первичных выражений (переменные и скобки)
        private Expr parsePrimary() {
            if (match(TokenType.LPAR)) {
                Expr expr = parseImpl();
                expect(TokenType.RPAR);
                return expr;
            }
            if (match(TokenType.VAR)) {
                return new VarExpr(previous().value);
            }
            throw new RuntimeException("Unexpected token: " + peek().type);
        }

        // Проверка текущего токена
        private boolean match(TokenType type) {
            if (check(type)) {
                pos++;
                return true;
            }
            return false;
        }

        // Проверка текущего токена без увеличения позиции
        private boolean check(TokenType type) {
            return !isAtEnd() && peek().type == type;
        }

        // Получение текущего токена
        private Token peek() {
            return tokens.get(pos);
        }

        // Получение предыдущего токена
        private Token previous() {
            return tokens.get(pos - 1);
        }

        // Проверка, достигнут ли конец списка токенов
        private boolean isAtEnd() {
            return pos >= tokens.size();
        }

        // Ожидание определённого токена
        private void expect(TokenType type) {
            if (!match(type)) {
                throw new RuntimeException("Expected token: " + type);
            }
        }
    }

    // Абстрактный класс для выражений
    abstract static class Expr {
        abstract String toPrefixString();
    }

    // Классы для различных типов выражений
    static class VarExpr extends Expr {
        final String name;

        VarExpr(String name) {
            this.name = name;
        }

        @Override
        String toPrefixString() {
            return name;
        }
    }

    static class ImplExpr extends Expr {
        final Expr left;
        final Expr right;

        ImplExpr(Expr left, Expr right) {
            this.left = left;
            this.right = right;
        }

        @Override
        String toPrefixString() {
            return "(->," + left.toPrefixString() + "," + right.toPrefixString() + ")";
        }
    }

    static class OrExpr extends Expr {
        final Expr left;
        final Expr right;

        OrExpr(Expr left, Expr right) {
            this.left = left;
            this.right = right;
        }

        @Override
        String toPrefixString() {
            return "(|," + left.toPrefixString() + "," + right.toPrefixString() + ")";
        }
    }

    static class AndExpr extends Expr {
        final Expr left;
        final Expr right;

        AndExpr(Expr left, Expr right) {
            this.left = left;
            this.right = right;
        }

        @Override
        String toPrefixString() {
            return "(&," + left.toPrefixString() + "," + right.toPrefixString() + ")";
        }
    }

    static class NotExpr extends Expr {
        final Expr expr;

        NotExpr(Expr expr) {
            this.expr = expr;
        }

        @Override
        String toPrefixString() {
            return "(!" + expr.toPrefixString() + ")";
        }
    }

    // Основная функция
    public static void main(String[] args) {
        Scanner scan = new Scanner(System.in);
        String input = scan.nextLine();
        List<Token> tokens = tokenize(input);
        Parser parser = new Parser(tokens);
        String result = parser.parse();
        System.out.println(result);
        scan.close();
    }
}










/* 
import java.util.*;

public class Parser {
    static Map<Character, Byte> PRIORITIES = new HashMap<>() {{
        put('&', (byte) 3);
        put('|', (byte) 4);
        put('-', (byte) 5);
    }};
    private int input;
    private char c;
    private static String expression;

    public static void main(String[] args) {
        Scanner scan = new Scanner(System.in);
        expression = scan.nextLine();
        scan.close();
        Parser parser = new Parser();
        String a = parser.parsing().replace("-", "->");
        System.out.println(a);        
    }

    public String parsing() {
        int end = expression.length();
        List<Byte> stackOfOperationPriority = new ArrayList<>();
        List<Character> stackOfOperationSign = new ArrayList<>();
        List<String> stackOfValues = new ArrayList<>();

        while (input < end) {
            c = expression.charAt(input++);

            if (Character.isWhitespace(c) || c == '\r' || c == '\n') {
                continue;
            }else if (Character.isAlphabetic(c)){
                input--;
                stackOfValues.add(isSomeWord(c));
                if (input >= end) {
                    break;
                } else {
                    c = expression.charAt(input++);
                }
            }
            else if ((c == '!')) {
                if (expression.charAt(input) == '(') {
                    input++;
                    String a = parsing();
                    stackOfValues.add("(!" + a + ")");
                } else {
                    stackOfValues.add("(!" + isSomeWord(expression.charAt(input)) + ")");
                    if (input < end - 1) {
                        c = expression.charAt(input++);
                    }
                }
                if (input >= end && c != ')') {
                    break;
                }
            }
            if (PRIORITIES.containsKey(c)) {
                while (!stackOfOperationSign.isEmpty()) {
                    byte priority = stackOfOperationPriority.get(stackOfOperationSign.size() - 1);
                    if (priority <= PRIORITIES.get(c) && priority != 5) {
                        String operation2 = stackOfValues.remove(stackOfValues.size() - 1);
                        String operation1 = stackOfValues.remove(stackOfValues.size() - 1);
                        char sign = stackOfOperationSign.remove(stackOfOperationSign.size() - 1);
                        stackOfOperationPriority.remove(stackOfOperationPriority.size() - 1);
                        stackOfValues.add("(" + sign + "," + operation1 + "," + operation2 + ")");
                    } else {
                        break;
                    }
                }
                stackOfOperationPriority.add(PRIORITIES.get(c));
                stackOfOperationSign.add(c);
                if (c == '-') {
                    input++;
                }
            }

            else if (c == ')') {
                while (stackOfValues.size() > 1) {
                    String operation2 = stackOfValues.remove(stackOfValues.size() - 1);
                    String operation1 = stackOfValues.remove(stackOfValues.size() - 1);
                    char sn = stackOfOperationSign.remove(stackOfOperationSign.size() - 1);
                    stackOfValues.add("(" + sn + "," + operation1 + "," + operation2 + ")");
                }
                return stackOfValues.remove(0);
            }
            else if (c == '(') {
                stackOfValues.add(parsing());
            }
        }
        while (stackOfValues.size() > 1) {
            String operation2 = stackOfValues.remove(stackOfValues.size() - 1);
            String operation1 = stackOfValues.remove(stackOfValues.size() - 1);
            char sn = stackOfOperationSign.remove(stackOfOperationSign.size()-1);
            stackOfValues.add("(" + sn + "," + operation1 + "," + operation2 + ")");
        }
        return stackOfValues.remove(0);
    }

    public String isSomeWord(char c) {
        StringBuilder sb = new StringBuilder();
        StringBuilder par = new StringBuilder();

        while (input < expression.length() && expression.charAt(input) != ')' && !Character.isWhitespace(expression.charAt(input)) && !PRIORITIES.containsKey(expression.charAt(input))) {
            if (expression.charAt(input) == '!') {
                par.append(')');
                sb.append("(!");
                input++;
            } else if (expression.charAt(input) == '(') {
                input++;
                String a = parsing();
                sb.append(a);
                // c = expression.charAt(input++);
            } else {
                sb.append(expression.charAt(input++));
            }
        }
        sb.append(par);
        return sb.toString();
    }
}

 */
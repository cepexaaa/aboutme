import java.util.*;

public class ProofParser {

    enum TokenType {
        VAR, DISJ, CONJ, IMPL, INVERSE, LPAR, RPAR, COMMA, TURN, HALT
    }

    static class Token {
        TokenType type;
        String value;

        Token(TokenType type, String value) {
            this.type = type;
            this.value = value;
        }
    }

    static List<Token> tokenize(String input) {
        List<Token> tokens = new ArrayList<>();
        int i = 0;
        while (i < input.length()) {
            char c = input.charAt(i);
            if (Character.isWhitespace(c) || c == '\t' || c == '\r') {
                i++;
            } else if (c == '|' && i + 1 < input.length() && input.charAt(i + 1) == '-') {
                tokens.add(new Token(TokenType.TURN, "|-"));
                i += 2;
            } else if (c == '|') {
                tokens.add(new Token(TokenType.DISJ, "|"));
                i++;
            } else if (c == '&') {
                tokens.add(new Token(TokenType.CONJ, "&"));
                i++;
            } else if (c == '-' && i + 1 < input.length() && input.charAt(i + 1) == '>') {
                tokens.add(new Token(TokenType.IMPL, "->"));
                i += 2;
            } else if (c == '!') {
                tokens.add(new Token(TokenType.INVERSE, "!"));
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
                tokens.add(new Token(TokenType.HALT, ""));
                i++;
            }
        }
        return tokens;
    }

    static class Parser {
        List<Token> tokens;
        int pos;

        Parser(List<Token> tokens) {
            this.tokens = tokens;
            this.pos = 0;
        }

        Token peek() {
            if (pos >= tokens.size()) {
                return new Token(TokenType.HALT, "");
            }
            return tokens.get(pos);
        }

        void consume() {
            if (pos < tokens.size()) {
                pos++;
            }
        }

        Expr parseE() {
            Expr d = parseD();
            if (peek().type == TokenType.IMPL) {
                consume();
                Expr e = parseE();
                return new ImplExpr(d, e);
            }
            return d;
        }

        Expr parseD() {
            Expr c = parseC();
            return parseDPrime(c);
        }

        Expr parseDPrime(Expr acc) {
            if (peek().type == TokenType.DISJ) {
                consume();
                Expr c = parseC();
                return parseDPrime(new DisjExpr(acc, c));
            }
            return acc;
        }

        Expr parseC() {
            Expr i = parseI();
            return parseCPrime(i);
        }

        Expr parseCPrime(Expr acc) {
            if (peek().type == TokenType.CONJ) {
                consume();
                Expr i = parseI();
                return parseCPrime(new ConjExpr(acc, i));
            }
            return acc;
        }

        Expr parseI() {
            if (peek().type == TokenType.INVERSE) {
                consume();
                Expr e = parseI();
                return new InverseExpr(e);
            } else if (peek().type == TokenType.LPAR) {
                consume();
                Expr e = parseE();
                if (peek().type != TokenType.RPAR) {
                    throw new RuntimeException("Expected ')'");
                }
                consume();
                return e;
            } else if (peek().type == TokenType.VAR) {
                Token t = peek();
                consume();
                return new VarExpr(t.value);
            } else {
                throw new RuntimeException("Unexpected token: " + peek().type);
            }
        }
    }

    abstract static class Expr {
        @Override
        abstract public String toString();

        // Проверка на равенство выражений
        @Override
        public boolean equals(Object obj) {
            if (obj instanceof Expr) {
                return this.toString().equals(obj.toString());
            }
            return false;
        }

        @Override
        public int hashCode() {
            return this.toString().hashCode();
        }
    }

    static class VarExpr extends Expr {
        String name;

        VarExpr(String name) {
            this.name = name;
        }

        @Override
        public String toString() {
            return name;
        }
    }

    static class DisjExpr extends Expr {
        Expr left, right;

        DisjExpr(Expr left, Expr right) {
            this.left = left;
            this.right = right;
        }

        @Override
        public String toString() {
            return "(" + left.toString() + "|" + right.toString() + ")";
        }
    }

    static class ConjExpr extends Expr {
        Expr left, right;

        ConjExpr(Expr left, Expr right) {
            this.left = left;
            this.right = right;
        }

        @Override
        public String toString() {
            return "(" + left.toString() + "&" + right.toString() + ")";
        }
    }

    static class ImplExpr extends Expr {
        Expr left, right;

        ImplExpr(Expr left, Expr right) {
            this.left = left;
            this.right = right;
        }

        @Override
        public String toString() {
            return "(" + left.toString() + "->" + right.toString() + ")";
        }
    }

    static class InverseExpr extends Expr {
        Expr expr;

        InverseExpr(Expr expr) {
            this.expr = expr;
        }

        @Override
        public String toString() {
            return "(!" + expr.toString() + ")";
        }
    }

    static class ProofLine {
        List<Expr> context;
        Expr expression;
        int lineNumber;
        String annotation;

        ProofLine(List<Expr> context, Expr expression, int lineNumber) {
            this.context = context;
            this.expression = expression;
            this.lineNumber = lineNumber;
            this.annotation = "";
        }

        @Override
        public String toString() {
            StringBuilder contextStr = new StringBuilder();
            for (Expr expr : context) {
                contextStr.append(expr.toString()).append(",");
            }
            if (contextStr.length() > 0) {
                contextStr.deleteCharAt(contextStr.length() - 1); // Удаляем последнюю запятую
            }
            return "[" + lineNumber + "] " + contextStr + "|-" + expression + " " + annotation;
        }
    }

    public static boolean compareContext(List<Expr> left, List<Expr> right) {
        if (left.size() != right.size()) {
            return false; 
        }
        for (Expr element : left) {
            if (Collections.frequency(left, element) != Collections.frequency(right, element)) {
                return false; 
            }
        }
        return true; 
    }

    static int checkAxiom(Expr expr) {
    if (expr instanceof ImplExpr) {
        ImplExpr impl1 = (ImplExpr) expr;
        if (impl1.right instanceof ImplExpr) {
            ImplExpr impl2 = (ImplExpr) impl1.right;
            if (impl1.left.equals(impl2.right)) {
                return 1;
            }
        }
    }
   // Проверка на вторую схему аксиом: (A -> B) -> (A -> B -> C) -> (A -> C)
if (expr instanceof ImplExpr) {
    ImplExpr impl1 = (ImplExpr) expr;
    if (impl1.left instanceof ImplExpr && impl1.right instanceof ImplExpr) {
        ImplExpr leftImpl = (ImplExpr) impl1.left;  // (A -> B)
        ImplExpr rightImpl = (ImplExpr) impl1.right; // (A -> B -> C) -> (A -> C)

        if (rightImpl.left instanceof ImplExpr && rightImpl.right instanceof ImplExpr) {
            ImplExpr rightLeftImpl = (ImplExpr) rightImpl.left;  // (A -> B -> C)
            ImplExpr rightRightImpl = (ImplExpr) rightImpl.right; // (A -> C)
            if (leftImpl.left.equals(rightLeftImpl.left) && // A == A
                rightLeftImpl.right instanceof ImplExpr &&  // (B -> C)
                ((ImplExpr) rightLeftImpl.right).left.equals(leftImpl.right) && // B == B
                rightRightImpl.left.equals(leftImpl.left) && // A == A
                rightRightImpl.right.equals(((ImplExpr) rightLeftImpl.right).right)) { // C == C
                return 2;
            }
        }
    }
}

    if (expr instanceof ImplExpr) {
        ImplExpr impl1 = (ImplExpr) expr;
        if (impl1.right instanceof ImplExpr) {
            ImplExpr impl2 = (ImplExpr) impl1.right;
            if (impl2.right instanceof ConjExpr) {
                ConjExpr conj = (ConjExpr) impl2.right;
                if (impl1.left.equals(conj.left) && impl2.left.equals(conj.right)) {
                    return 3;
                }
            }
        }
    }

    if (expr instanceof ImplExpr) {
        ImplExpr impl = (ImplExpr) expr;
        if (impl.left instanceof ConjExpr) {
            ConjExpr conj = (ConjExpr) impl.left;
            if (conj.left.equals(impl.right)) {
                return 4;
            }
        }
    }

    if (expr instanceof ImplExpr) {
        ImplExpr impl = (ImplExpr) expr;
        if (impl.left instanceof ConjExpr) {
            ConjExpr conj = (ConjExpr) impl.left;
            if (conj.right.equals(impl.right)) {
                return 5;
            }
        }
    }

    if (expr instanceof ImplExpr) {
        ImplExpr impl = (ImplExpr) expr;
        if (impl.right instanceof DisjExpr) {
            DisjExpr disj = (DisjExpr) impl.right;
            if (impl.left.equals(disj.left)) {
                return 6;
            }
        }
    }

    if (expr instanceof ImplExpr) {
        ImplExpr impl = (ImplExpr) expr;
        if (impl.right instanceof DisjExpr) {
            DisjExpr disj = (DisjExpr) impl.right;
            if (impl.left.equals(disj.right)) {
                return 7;
            }
        }
    }

    if (expr instanceof ImplExpr) {
        ImplExpr impl1 = (ImplExpr) expr;
        if (impl1.left instanceof ImplExpr && impl1.right instanceof ImplExpr) {
            ImplExpr leftImpl = (ImplExpr) impl1.left;
            ImplExpr rightImpl = (ImplExpr) impl1.right;
            if (rightImpl.left instanceof ImplExpr && rightImpl.right instanceof ImplExpr) {
                ImplExpr rightLeftImpl = (ImplExpr) rightImpl.left;
                ImplExpr rightRightImpl = (ImplExpr) rightImpl.right;
                if (rightRightImpl.left instanceof DisjExpr) {
                    DisjExpr disj = (DisjExpr) rightRightImpl.left;
                    if (leftImpl.right.equals(rightLeftImpl.right) &&
                        rightLeftImpl.right.equals(rightRightImpl.right) &&
                        leftImpl.left.equals(disj.left) &&
                        rightLeftImpl.left.equals(disj.right)) {
                        return 8;
                    }
                }
            }
        }
    }

    if (expr instanceof ImplExpr) {
        ImplExpr impl1 = (ImplExpr) expr;
        if (impl1.left instanceof ImplExpr && impl1.right instanceof ImplExpr) {
            ImplExpr leftImpl = (ImplExpr) impl1.left;
            ImplExpr rightImpl = (ImplExpr) impl1.right;
            if (rightImpl.left instanceof ImplExpr && rightImpl.right instanceof InverseExpr) {
                ImplExpr rightLeftImpl = (ImplExpr) rightImpl.left;
                InverseExpr inverse = (InverseExpr) rightImpl.right;
                if (leftImpl.left.equals(rightLeftImpl.left) &&
                    leftImpl.right.equals(((InverseExpr) rightLeftImpl.right).expr) &&
                    inverse.expr.equals(leftImpl.left)) {
                    return 9;
                }
            }
        }
    }

    if (expr instanceof ImplExpr) {
        ImplExpr impl = (ImplExpr) expr;
        if (impl.left instanceof InverseExpr) {
            InverseExpr inverse1 = (InverseExpr) impl.left;
            if (inverse1.expr instanceof InverseExpr) {
                InverseExpr inverse2 = (InverseExpr) inverse1.expr;
                if (inverse2.expr.equals(impl.right)) {
                    return 10;
                }
            }
        }
    }

    return -1; // Если не аксиома
}

    static boolean isHypothesis(ProofLine line) {
        return line.context.contains(line.expression);
    }


    static String checkModusPonens(List<ProofLine> proof, ProofLine line) {
        ProofLine bestImplLine = null; // Строка с импликацией
        ProofLine bestPremiseLine = null; // Строка с посылкой

        for (ProofLine prevLine : proof) {
            if (prevLine.expression instanceof ImplExpr) {
                ImplExpr impl = (ImplExpr) prevLine.expression;
                if (impl.right.equals(line.expression)) {
                    for (ProofLine prevLine2 : proof) {
                        if (prevLine2.expression.equals(impl.left)) {
                            if (prevLine.lineNumber >= line.lineNumber || prevLine2.lineNumber >= line.lineNumber) {
                                break;
                            }
                            // Проверяем, что контексты совпадают
                            if (compareContext(prevLine.context, line.context) && compareContext(prevLine2.context, line.context)) {
                                // Если обе строки корректны, сразу возвращаем результат
                                if (!prevLine.annotation.equals("[Incorrect]") && !prevLine2.annotation.equals("[Incorrect]")) {
                                    return "[M.P. " + prevLine2.lineNumber + ", " + prevLine.lineNumber + "]";
                                }
                                if (bestImplLine == null || bestPremiseLine == null) {
                                    bestImplLine = prevLine;
                                    bestPremiseLine = prevLine2;
                                }
                            }
                        }
                    }
                }
            }
        }
        if (bestImplLine != null && bestPremiseLine != null) {
            return "[M.P. " + bestPremiseLine.lineNumber + ", " + bestImplLine.lineNumber + "; from Incorrect]";
        }
        return "[Incorrect]";
    }
    static String checkDeduction(List<ProofLine> proof, ProofLine line) {
        // Получаем максимальный контекст и оставшееся выражение для текущей строки
        Pair<List<Expr>, Expr> currentResult = getMaxContextAndExpression(line.expression, line.context);
        List<Expr> maxContext = currentResult.getLeft();
        Expr remainingExpr = currentResult.getRight();

        // Ищем строку, которая может быть использована для дедукции
        ProofLine bestMatch = null;
        for (ProofLine prevLine : proof) {
            if (prevLine.lineNumber >= line.lineNumber) {
                break;
            }
            Pair<List<Expr>, Expr> prevResult = getMaxContextAndExpression(prevLine.expression, prevLine.context);
            List<Expr> prevMaxContext = prevResult.getLeft();
            Expr prevRemainingExpr = prevResult.getRight();

            if (compareContext(prevMaxContext, maxContext)) {
                if (prevRemainingExpr.equals(remainingExpr)) {
                    if (!prevLine.annotation.equals("[Incorrect]")) {
                        return "[Ded. " + prevLine.lineNumber + "]";
                    }
                    if (bestMatch == null) {
                        bestMatch = prevLine;
                    }
                }
            }
        }
        if (bestMatch != null) {
            return "[Ded. " + bestMatch.lineNumber + "; from Incorrect]";
        }

        return "[Incorrect]";
    }

    static Pair<List<Expr>, Expr> getMaxContextAndExpression(Expr expression, List<Expr> context) {
    List<Expr> maxContext = new ArrayList<>(context);

    if (expression instanceof ImplExpr) {
        ImplExpr impl = (ImplExpr) expression;
        maxContext.add(impl.left);
        return getMaxContextAndExpression(impl.right, maxContext);
    }

    return new Pair<>(maxContext, expression);
}

    static class Pair<L, R> {
    private final L left;
    private final R right;

    Pair(L left, R right) {
        this.left = left;
        this.right = right;
    }

    public L getLeft() {
        return left;
    }

    public R getRight() {
        return right;
    }
}

    static List<Expr> getMaxContext(Expr expression, List<Expr> context) {
    List<Expr> maxContext = new ArrayList<>(context);

    if (expression instanceof ImplExpr) {
        ImplExpr impl = (ImplExpr) expression;
        maxContext.add(impl.left);
        return getMaxContext(impl.right, maxContext);
    }

    return maxContext;
}

    static void annotateProof(List<ProofLine> proof) {
        for (ProofLine line : proof) {
            int axiomNum = checkAxiom(line.expression);
            if (axiomNum != -1) {
                line.annotation = "[Ax. sch. " + axiomNum + "]";
            } else if (isHypothesis(line)) {
                line.annotation = "[Hyp. " + (line.context.indexOf(line.expression) + 1) + "]";
            } else {
                String mpAnnotation = checkModusPonens(proof, line);
                if (!mpAnnotation.equals("[Incorrect]")) {
                    line.annotation = mpAnnotation;
                } else {
                    String dedAnnotation = checkDeduction(proof, line);
                    line.annotation = dedAnnotation;
                }
            }
        }
    }

    public static void main(String[] args) {
        List<ProofLine> proof = new ArrayList<>();
        int lineNumber = 1;
        Scanner scanner = new Scanner(System.in);    
        while (scanner.hasNextLine()) {
            String input = scanner.nextLine().trim();     
            if (input.isEmpty()) {
                continue;
            }    
            List<Token> tokens = tokenize(input);
            Parser parser = new Parser(tokens);
            List<Expr> context = new ArrayList<>();
            while (parser.peek().type != TokenType.TURN) {
                context.add(parser.parseE());
                if (parser.peek().type == TokenType.COMMA) {
                    parser.consume();
                }
            }
            parser.consume(); // Пропускаем "|-"
            Expr expression = parser.parseE();
    
            proof.add(new ProofLine(context, expression, lineNumber++));
        }
        scanner.close();
        annotateProof(proof);
        for (ProofLine line : proof) {
            System.out.println(line);
        }
    }
}
//если МР выдаёт как следствие из неверного, то надо ли проверять дедукцию, которая может быть из корректного?

/* |- (A -> B) -> (A -> B -> C) -> (A -> C) 
   |- A->B->A&B
   |- A&B->A
   |- A&B->B
   |- A->A|B
   |- B->A|B
   |- (A -> C) -> (B -> C) -> (A | B -> C)
   |- (A -> B) -> (A -> !B) -> !A
   |- !!A -> A 

  
  
/* // Проверка на Modus Ponens
static String checkModusPonens(List<ProofLine> proof, ProofLine line) {
    for (ProofLine prevLine : proof) {
        if (prevLine.expression instanceof ImplExpr) {
            ImplExpr impl = (ImplExpr) prevLine.expression;
            if (impl.right.equals(line.expression)) {
                for (ProofLine prevLine2 : proof) {
                    if (prevLine2.expression.equals(impl.left)) {
                        // Проверяем, что контексты совпадают
                        if (compareContext(prevLine.context, line.context) && compareContext(prevLine2.context, line.context) ) {
                            // Проверяем, являются ли обе строки корректными
                            boolean isPrevLineCorrect = !prevLine.annotation.equals("[Incorrect]");
                            boolean isPrevLine2Correct = !prevLine2.annotation.equals("[Incorrect]");

                            if (isPrevLineCorrect && isPrevLine2Correct) {
                                return "[M.P. " + prevLine2.lineNumber + ", " + prevLine.lineNumber + "]";
                            } else {
                                return "[M.P. " + prevLine2.lineNumber + ", " + prevLine.lineNumber + "; from Incorrect]";
                            }
                        }
                    }
                }
            }
        }
    }
    return "[Incorrect]";
} */

//if (compareContext(prevMaxContext, maxContext)) {


// Проверка на теорему о дедукции
 
// Проверка на теорему о дедукции
/* static String checkDeduction(List<ProofLine> proof, ProofLine line) {
    // Получаем максимальный контекст и оставшееся выражение для текущей строки
    Pair<List<Expr>, Expr> currentResult = getMaxContextAndExpression(line.expression, line.context);
    List<Expr> maxContext = currentResult.getLeft();
    Expr remainingExpr = currentResult.getRight();

    // Ищем строку, которая может быть использована для дедукции
    for (ProofLine prevLine : proof) {
        if (prevLine.lineNumber == line.lineNumber) {
            return "[Incorrect]";
        }
        // Получаем максимальный контекст и оставшееся выражение для предыдущей строки
        Pair<List<Expr>, Expr> prevResult = getMaxContextAndExpression(prevLine.expression, prevLine.context);
        List<Expr> prevMaxContext = prevResult.getLeft();
        Expr prevRemainingExpr = prevResult.getRight();

        // Проверяем, что контексты совпадают
        if (compareContext(prevMaxContext, maxContext)) {
            // Проверяем, что оставшиеся выражения совпадают
            if (prevRemainingExpr.equals(remainingExpr)) {
                // Проверяем, является ли строка корректной
                boolean isPrevLineCorrect = !prevLine.annotation.equals("[Incorrect]");

                if (isPrevLineCorrect) {
                    return "[Ded. " + prevLine.lineNumber + "]";
                } else {
                    return "[Ded. " + prevLine.lineNumber + "; from Incorrect]";
                }
            }
        }
    }

    return "[Incorrect]";
} */
  
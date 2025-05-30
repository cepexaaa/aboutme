package expression.generic;

import expression.generic.TypesT.*;

public class GenericTabulator implements Tabulator {

    @Override
    public Object[][][] tabulate(String mode, String expression, int x1, int x2, int y1, int y2, int z1, int z2){
        InterfaceT<?> type;
        switch (mode) {
            case "i" -> type = new IntegerT(true);
            case "d" -> type = new DoubleT();
            case "bi" -> type = new BigIntegerT();
            case "u" -> type = new IntegerT(false);
            case "b" -> type = new ByteT();
            default -> type = new BooleanT();
        }
        return table(type, expression, x1, x2, y1, y2, z1, z2);
    }
    private <T> Object[][][] table(InterfaceT<T> type, String expression, int x1, int x2, int y1, int y2, int z1, int z2) {
        ExpressionParser<T> parser = new ExpressionParser<>(type);
        SomeExpressions<T> expression_p = parser.parse(expression);
        Object[][][] result = new Object[x2 - x1 + 1][y2 - y1 + 1][z2 - z1 + 1];
        for (int i = 0; i < x2 - x1 + 1; i++) {
            for (int j = 0; j < y2 - y1 + 1; j++) {
                for (int k = 0; k < z2 - z1 + 1; k++) {
                    try {
                        result[i][j][k] = expression_p.evaluate(type.valueOf((x1 + i)), type.valueOf((y1 + j)),
                                type.valueOf((z1 + k)));
                    } catch (Exception e) {
                        result[i][j][k] = null;
                    }

                }
            }
        }
        return result;
    }
}





























//    public static InterfaceT<?> getType(String x) {
//        return switch (x) {
//            case "i" -> new IntegerT(true);
//            case "d" -> new DoubleT();
//            case "bi" -> new BigIntegerT();
//            case "u" -> new IntegerT(false);
//            case "b" -> new ByteT();
//            case "bool" -> new BooleanT();
//            default -> null;
//        };
//    }

//return table(KindOfType.getType(mode), expression, x1, x2, y1, y2, z1, z2);

//    private <T> Object[][][] table(InterfaceT<T> type, String expression, int x1, int x2, int y1, int y2, int z1, int z2) {
//        ExpressionParser<T> parser = new ExpressionParser<>(type);
//        SomeExpressions<T> expression_p = parser.parse(expression);
//        Object[][][] result = new Object[x2 - x1 + 1][y2 - y1 + 1][z2 - z1 + 1];
//        for (int i = 0; i < x2 - x1 + 1; i++) {
//            for (int j = 0; j < y2 - y1 + 1; j++) {
//                for (int k = 0; k < z2 - z1 + 1; k++) {
//                    result[i][j][k] = expression_p.evaluate(type.valueOf((x1 + i)), type.valueOf((y1 + j)),
//                            type.valueOf((z1 + k)));
//                }
//            }
//        }
//        return result;
//    }

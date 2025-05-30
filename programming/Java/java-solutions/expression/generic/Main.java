package expression.generic;

public class Main {
    public static void main(String[] args) throws Exception {
        String mode = args[0].substring(1);
        String expression = args[1];
        GenericTabulator tabulator = new GenericTabulator();
        tabulator.tabulate(mode, expression, -2, 2, -2, 2, -2, 2);
        for(int i = 0; i < 5; i++) {
            for (int j =0; j < 5; j++) {
                for (int k = 0; k < 5; k++) {
                    System.out.println(tabulator.tabulate(mode, expression, -2, 2, -2, 2, -2, 2)[i][j][k]);
                }
            }
        }
    }
}

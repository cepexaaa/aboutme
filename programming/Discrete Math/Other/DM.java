package Other;

public class DM {
    public static void main(String[] args) {
        double H = 1.2750;
        double eps = 0.0000001;
        int n = 400;
        for (int a = 0; a < n; a++) {
            for (int b = a + 1; b < n; b++) {
                for (int c = b + 1; c < n; c++) {
                    for (int d = c + 1; d < n; d++) {
                        int sum = a + b + c + d;
                        double fa = (double) a/sum;
                        double fb = (double) b/sum;
                        double fc = (double) c/sum;
                        double fd = (double) d/sum;
                        double h = -(fa*log2(fa) + fb*log2(fb) + fc*log2(fc) + fd*log2(fd));
                        if (Math.abs(H - h) < eps){
                            System.out.println(h);
                            System.out.println(a + " " + b + " " + c + " " + d);
                        }
                    }
                }
            }
        }
    }

    private static double log2(double fa) {
        return Math.log(fa)/Math.log(2);
    }
}//0 9 165 329



































// public class  {
//    public static void main(String[] args) {
//        double H = 1.5890;
//        double eps = 0.00001;
//        int n = 333;
//
//        for (int a = 1; a < n; a++) {
//            for (int b = a + 1; b < n; b++) {
//                for (int c = b + 1; c < n; c++) {
//                    for (int d = c + 1; d < n; d++) {
//                        int sm = a + b + c + d;
//                        double fa = (double) a / sm;
//                        double fb = (double) b / sm;
//                        double fc = (double) c / sm;
//                        double fd = (double) d / sm;
//
//                        double h = -fa * log2(fa) - fb * log2(fb) - fc * log2(fc) - fd * log2(fd);
//
//                        if (Math.abs(H - h) < eps) {
//                            System.out.println(a + " " + b + " " + c + " " + d);
//                            System.exit(0);
//                        }
//                    }
//                }
//            }
//        }
//    }
//
//    private static double log2(double x) {
//        return Math.log(x) / Math.log(2);
//    }
//}

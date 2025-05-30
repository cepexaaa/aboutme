package antropiya.src;//import java.util.Scanner;
//
//public class Main {
//    public static void main(String[] args) {
//        double sum = 0;
//        Scanner scan = new Scanner(System.in);
//        int n = 26;//scan.nextInt();
//        for (int i = 0; i < n; i++) {
//            double p = scan.nextDouble();
//            sum += (p*Math.log(p)/Math.log(2));
//        }
//        System.out.println(sum);
//    }
//}//0.0817 0.0149 0.0278 0.0425 0.0127 0.0223 0.0202 0.0609 0.0697 0.0015 0.0077 0.0403 0.0241 0.0675 0.0751 0.0193 0.001 0.0599 0.0633 0.0906 0.0276 0.0098 0.0236 0.0015 0.0197 0.0007

import java.util.Scanner;

public class Main {
    public static void main(String[] args) {
        double sum = 0;
        int[] a = new int[]{1, 4, 12, 14};
        int count = 0;
        for (int i = 0; i < a.length; i++) {
            count+=a[i];
        }
//        String input = "" + (29/(29 + 30 + 31 + 117)) + " " + (30/(29 + 30 + 31 + 117)) + " " + (31/(29 + 30 + 31 + 117)) + " " + (117/(29 + 30 + 31 + 117)) + " " ;
        //double[] input = new double[]{ (29.0/(29 + 30 + 31 + 117)), (30.0/(29 + 30 + 31 + 117)), (31.0/(29 + 30 + 31 + 117)), (117.0/(29 + 30 + 31 + 117))};
        double[] input = new double[a.length];
        for (int i = 0; i < a.length; i++) {
            input[i] = (double)a[i]/count;
        }
//        Scanner scan = new Scanner(input);
//
//        while (scan.hasNext()) {
//            double p = Double.parseDouble(scan.next());
//            sum += p * Math.log(p) / Math.log(2);
//        }
        for (int i = 0; i < input.length; i++) {
            sum += input[i] * Math.log(input[i]) / Math.log(2);
        }
//        for (int i = 0; i < 26; i++) {
//            sum += 0.038 * Math.log(0.038) / Math.log(2);
//        }
        System.out.println(sum);
    }
}//-3.8781400790744147
//-4.661242489963879

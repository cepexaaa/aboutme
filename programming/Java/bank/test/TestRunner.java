package info.kgeorgiy.ja.kubesh.bank.test;

import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;

public class TestRunner {
    public static void main(String[] args) {
        Class[] testClasses = {
                BankTest.class,
                PersonTest.class,
                RemoteTest.class,
                LocalTest.class,
                MoveAmountTest.class
        };

        long startTime = System.currentTimeMillis();

        Result result = JUnitCore.runClasses(testClasses);

        for (Failure failure : result.getFailures()) {
            System.out.println(failure.toString());
            System.out.println();
        }

        long endTime = System.currentTimeMillis();
        long elapsedTime = endTime - startTime;

        System.out.println("\n------------ Results ------------");
        System.out.println("Elapsed time: " + elapsedTime + " ms");
        System.out.println("Tests run: " + result.getRunCount());
        System.out.println("Tests failed: " + result.getFailureCount());
        System.out.println("Tests successful: " + (result.getRunCount() - result.getFailureCount()));

        System.exit(result.wasSuccessful() ? 0 : 1);
    }
}

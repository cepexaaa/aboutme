package info.kgeorgiy.ja.kubesh.bank.test;

import info.kgeorgiy.ja.kubesh.bank.rmi.Account;
import info.kgeorgiy.ja.kubesh.bank.rmi.Bank;
import info.kgeorgiy.ja.kubesh.bank.rmi.Person;
import info.kgeorgiy.ja.kubesh.bank.rmi.RemoteBank;
import org.junit.Before;
import org.junit.Test;

import java.rmi.RemoteException;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

public class MoveAmountTest {
    private Bank bank;

    @Before
    public void setUp() throws RemoteException {
        bank = new RemoteBank(8888);
    }

    @Test
    public void testMoveAmount() throws RemoteException {
        Person person = bank.createPerson("John", "Doe", "1234567890");
        Person person2 = bank.createPerson("John2", "Doe2", "111");
        Account acc1 = person.createAccount("acc1");
                acc1.setAmount(100);
        Account acc2 = person2.createAccount("acc2");
                acc2.setAmount(1000);
        System.out.println("Person 1 - money: " + acc1.getAmount());
        System.out.println("Person 2 - money: " + acc2.getAmount());
        person2.moveAmount("John", "Doe", 100);
        System.out.println("Sending money...");
        System.out.println("Person 1 - money: " + person.getAccount("acc1").getAmount());
        System.out.println("Person 2 - money: " + person2.getAccount("acc2").getAmount());

        assertEquals("Accounts should be synchronized", 200, person.getAccount("acc1").getAmount());
        assertEquals("Accounts should be synchronized", 900, person2.getAccount("acc2").getAmount());
        System.out.println("Money sent successfully");
    }

    @Test
    public void testMoveAmountToPersonWithoutAccounts() throws RemoteException {
        Person sender = bank.createPerson("John", "Doe", "sender");
        Person receiver = bank.createPerson("Alice", "Smith", "receiver");
        sender.createAccount("acc1").setAmount(100);

        try {
            sender.moveAmount("Alice", "Smith", 50);
            fail("Should throw exception when moving money to person without accounts");
        } catch (RuntimeException e) {
            // Expected
        }

        assertEquals(100, sender.getAccount("acc1").getAmount());
    }

    @Test
    public void testMoveAmountWithInsufficientFunds() throws RemoteException {
        Person sender = bank.createPerson("John", "Doe", "sender");
        Person receiver = bank.createPerson("Alice", "Smith", "receiver");
        sender.createAccount("acc1").setAmount(100);
        receiver.createAccount("acc1").setAmount(0);

        sender.moveAmount("Alice", "Smith", 150);

        assertEquals(-50, sender.getAccount("acc1").getAmount());
        assertEquals(150, receiver.getAccount("acc1").getAmount());
    }

    @Test
    public void testConcurrentMoneyTransfers() throws RemoteException, InterruptedException {
        Person person1 = bank.createPerson("John", "Doe", "p1");
        Person person2 = bank.createPerson("Alice", "Smith", "p2");
        person1.createAccount("acc").setAmount(1000);
        person2.createAccount("acc").setAmount(1000);

        final int threadCount = 10;
        Thread[] threads = new Thread[threadCount];

        for (int i = 0; i < threadCount; i++) {
            threads[i] = new Thread(() -> {
                try {
                    Person p1 = bank.getRemotePerson("p1");
                    Person p2 = bank.getRemotePerson("p2");
                    p1.moveAmount("Alice", "Smith", 10);
                } catch (RemoteException e) {
                    fail("RemoteException during money transfer");
                }
            });
            threads[i].start();
        }

        for (Thread thread : threads) {
            thread.join();
        }

        assertEquals(1000 - 10 * threadCount, person1.getAccount("acc").getAmount());
        assertEquals(1000 + 10 * threadCount, person2.getAccount("acc").getAmount());
    }
}

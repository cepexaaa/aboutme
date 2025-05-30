package info.kgeorgiy.ja.kubesh.bank.test;

import info.kgeorgiy.ja.kubesh.bank.rmi.*;
import org.junit.Before;
import org.junit.Test;
import java.rmi.RemoteException;
import static org.junit.Assert.*;

public class RemoteTest {
    private Bank bank;
    private final int port = 8888;

    @Before
    public void setUp() throws RemoteException {
        bank = new RemoteBank(port);
    }

    @Test
    public void testRemotePersonAccountSync() throws RemoteException {
        Person person = bank.createPerson("John", "Doe", "123");
        Account account = person.createAccount("acc1");
        account.setAmount(100);

        Account sameAccount = person.getAccount("acc1");
        assertEquals("Accounts should be synchronized", 100, sameAccount.getAmount());
    }

    @Test
    public void testManyRemotePersons() throws RemoteException {
        Person person = bank.createPerson("John", "Doe", "123");
        Account account = person.createAccount("acc1");
        account.setAmount(100);
        Person anotherPerson = bank.getRemotePerson("123");
        assertEquals("Accounts should be synchronized", 100, anotherPerson.getAccount("acc1").getAmount());
        account.addAmount(10000);
        assertEquals("Accounts should be synchronized", 10100, anotherPerson.getAccount("acc1").getAmount());
    }

    @Test
    public void testAccountIdFormat() throws RemoteException {
        Person person = bank.createPerson("John", "Doe", "pass123");
        Account account = person.createAccount("acc1");
        assertEquals("Account ID should be in format personId:subId",
                "pass123:acc1", account.getId());
    }

    @Test
    public void testRemotePersonBankReference() throws RemoteException {
        Person person = bank.createPerson("John", "Doe", "123");
        Bank personBank = person.getBank();
        assertNotNull("Person should have reference to bank", personBank);
        assertEquals("Bank references should be equal", bank, personBank);
    }

    @Test
    public void testConcurrentRemotePersonAccess() throws RemoteException, InterruptedException {
        final Person person = bank.createPerson("John", "Doe", "123");
        person.createAccount("acc1").setAmount(1000);

        final int THREAD_COUNT = 10;
        Thread[] threads = new Thread[THREAD_COUNT];
        final int[] amounts = new int[THREAD_COUNT];

        for (int i = 0; i < THREAD_COUNT; i++) {
            threads[i] = new Thread(() -> {
                try {
                    Person p = bank.getRemotePerson("123");
                    p.getAccount("acc1").addAmount(10);
                } catch (RemoteException e) {
                    fail("RemoteException during test");
                }
            });
            threads[i].start();
        }

        for (Thread thread : threads) {
            thread.join();
        }

        assertEquals(1000 + 10 * THREAD_COUNT, person.getAccount("acc1").getAmount());
    }
}
package info.kgeorgiy.ja.kubesh.bank.test;

import info.kgeorgiy.ja.kubesh.bank.rmi.*;
import org.junit.Before;
import org.junit.Test;
import java.rmi.RemoteException;
import static org.junit.Assert.*;

public class BankTest {
    private Bank bank;
    private final int port = 8888;

    @Before
    public void setUp() throws RemoteException {
        bank = new RemoteBank(port);
    }

    @Test
    public void testCreateAccount() throws RemoteException {
        Account account = bank.createAccount("123");
        assertNotNull("Account should be created", account);
        assertEquals("Account ID should match", "123", account.getId());
    }

    @Test
    public void testGetNonExistingAccount() throws RemoteException {
        assertNull("Non-existing account should be null", bank.getAccount("nonexistent"));
    }

    @Test
    public void testCreateExistingAccount() throws RemoteException {
        Account account1 = bank.createAccount("123");
        Account account2 = bank.createAccount("123");
        assertSame("Should return same account for same ID", account1, account2);
    }

    @Test
    public void testConcurrentAccountCreation() throws RemoteException, InterruptedException {
        final int THREAD_COUNT = 10;
        final Account[] accounts = new Account[THREAD_COUNT];

        Thread[] threads = new Thread[THREAD_COUNT];
        for (int i = 0; i < THREAD_COUNT; i++) {
            final int index = i;
            threads[i] = new Thread(() -> {
                try {
                    accounts[index] = bank.createAccount("concurrent");
                } catch (RemoteException e) {
                    fail("RemoteException during account creation");
                }
            });
            threads[i].start();
        }

        for (Thread thread : threads) {
            thread.join();
        }

        for (int i = 1; i < THREAD_COUNT; i++) {
            assertSame("All threads should get same account instance", accounts[0], accounts[i]);
        }
    }
}
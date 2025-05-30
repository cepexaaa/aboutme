package info.kgeorgiy.ja.kubesh.bank.test;

import info.kgeorgiy.ja.kubesh.bank.rmi.*;
import org.junit.Before;
import org.junit.Test;
import java.rmi.RemoteException;
import static org.junit.Assert.*;

public class LocalTest {
    private Bank bank;
    private final int port = 8888;

    @Before
    public void setUp() throws RemoteException {
        bank = new RemoteBank(port);
    }

    @Test
    public void testLocalPersonIsolation() throws RemoteException {
        Person remote = bank.createPerson("John", "Doe", "123");
        remote.createAccount("acc1").setAmount(100);

        Person local = bank.getLocalPerson("123");
        local.getAccount("acc1").setAmount(200);

        assertEquals("Local change should be isolated",
                200, local.getAccount("acc1").getAmount());
        assertEquals("Remote should not see local changes",
                100, remote.getAccount("acc1").getAmount());
    }

    @Test
    public void testManyLocalPersonsIsolation() throws RemoteException {
        Person remote = bank.createPerson("John", "Doe", "123");
        remote.createAccount("acc1").setAmount(100);

        Person local1 = bank.getLocalPerson("123");
        Person local2 = bank.getLocalPerson("123");
        local1.getAccount("acc1").setAmount(200);
        local2.getAccount("acc1").setAmount(300);

        assertEquals("Local change should be isolated",200, local1.getAccount("acc1").getAmount());
        assertEquals("Local change should be isolated",300, local2.getAccount("acc1").getAmount());
    }

    @Test
    public void testLocalChanges() throws RemoteException {
        Person remote = bank.createPerson("John", "Doe", "123");
        Person local = bank.getLocalPerson("123");
        Account localAccount = local.createAccount("acc1");
        Account remoteAccount = remote.getAccount("acc1");

        assertEquals("Local change should be isolated",0, localAccount.getAmount());
        assertNull(remoteAccount);
    }

    @Test
    public void testRemoteChangesToLocal() throws RemoteException {
        Person remote = bank.createPerson("John", "Doe", "123");
        Person local = bank.getLocalPerson("123");
        Account remoteAccount = remote.createAccount("acc1");
        Account localAccount = local.getAccount("acc1");

        assertEquals("Local change should be isolated",0, remoteAccount.getAmount());
        assertNull(localAccount);
    }

    @Test
    public void testLocalPersonSnapshot() throws RemoteException {
        Person remote = bank.createPerson("John", "Doe", "123");
        remote.createAccount("acc1").setAmount(100);

        Person local1 = bank.getLocalPerson("123");
        remote.getAccount("acc1").setAmount(200);

        Person local2 = bank.getLocalPerson("123");

        assertEquals("First local copy should have initial amount",
                100, local1.getAccount("acc1").getAmount());
        assertEquals("Second local copy should have updated amount",
                200, local2.getAccount("acc1").getAmount());
    }

    @Test
    public void testLocalPersonSerialization() throws RemoteException {
        Person remote = bank.createPerson("John", "Doe", "123");
        remote.createAccount("acc1").setAmount(100);

        Person local = bank.getLocalPerson("123");
        local.getAccount("acc1").setAmount(200);

        LocalPerson serializedCopy = new LocalPerson(local);

        assertEquals("Serialized copy should preserve local changes",
                200, serializedCopy.getAccount("acc1").getAmount());
    }
}


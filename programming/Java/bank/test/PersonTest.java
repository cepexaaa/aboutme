package info.kgeorgiy.ja.kubesh.bank.test;

import info.kgeorgiy.ja.kubesh.bank.rmi.*;
import org.junit.Before;
import org.junit.Test;
import java.rmi.RemoteException;
import static org.junit.Assert.*;

public class PersonTest {
    private Bank bank;
    private final int port = 8888;

    @Before
    public void setUp() throws RemoteException {
        bank = new RemoteBank(port);
    }

    @Test
    public void testCreatePerson() throws RemoteException {
        Person person = bank.createPerson("John", "Doe", "1234567890");
        assertNotNull("Person should be created", person);
        assertEquals("Name should match", "John", person.getName());
        assertEquals("Surname should match", "Doe", person.getSurname());
    }

    @Test
    public void testGetLocalPerson() throws RemoteException {
        bank.createPerson("John", "Doe", "1234567890");
        Person local = bank.getLocalPerson("1234567890");
        assertNotNull("Local person should be created", local);
        assertEquals("Name should match", "John", local.getName());
    }

    @Test
    public void testCreatePersonWithExistingPassport() throws RemoteException {
        Person first = bank.createPerson("John", "Doe", "123");
        Person second = bank.createPerson("Different", "Name", "123");
        assertEquals(first, second);
    }

    @Test
    public void testGetNonExistingPerson() throws RemoteException {
        assertNull(bank.getRemotePerson("nonexistent"));
        assertNull(bank.getLocalPerson("nonexistent"));
    }

    @Test
    public void testPersonAccountsList() throws RemoteException {
        Person person = bank.createPerson("John", "Doe", "123");
        person.createAccount("acc1");
        person.createAccount("acc2");

        assertEquals(2, person.getAccounts().size());
        assertTrue(person.getAccounts().contains("acc1"));
        assertTrue(person.getAccounts().contains("acc2"));
    }
}
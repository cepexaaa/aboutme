package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.rmi.RemoteException;
import java.rmi.server.UnicastRemoteObject;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

public class RemoteBank implements Bank {
    private final int port;
    private final ConcurrentMap<String, Account> accounts = new ConcurrentHashMap<>();
    private final ConcurrentMap<String, Person> persons = new ConcurrentHashMap<>();

    /**
     * Creates new remote bank
     * @param port where the bank will send and get signals
     */
    public RemoteBank(final int port) {
        this.port = port;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Account createAccount(final String id) throws RemoteException {
        System.out.println("Creating account " + id);
        final Account account = new RemoteAccount(id);
        if (accounts.putIfAbsent(id, account) == null) {
            UnicastRemoteObject.exportObject(account, port);
            return account;
        } else {
            return getAccount(id);
        }
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Account getAccount(final String id) {
        System.out.println("Retrieving account " + id);
        return accounts.get(id);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Person getRemotePerson(String passportNumber) throws RemoteException {
        Person person = persons.get(passportNumber);
        if (person != null) {
            System.out.println("Retrieving remote person " + passportNumber);
        }
        return person;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Person getLocalPerson(String passportNumber) throws RemoteException {
        Person remotePerson = persons.get(passportNumber);
        if (remotePerson != null) {
            System.out.println("Creating local copy of person " + passportNumber);
            return new LocalPerson(remotePerson);
        }
        return null;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Person createPerson(String name, String surname, String passportNumber) throws RemoteException {
        System.out.println("Creating person " + passportNumber);
        Person person = new RemotePerson(name, surname, passportNumber, this, port);
        if (persons.putIfAbsent(passportNumber, person) == null) {
            return person;
        } else {
            return getRemotePerson(passportNumber);
        }
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Person getPersonByNameAndSurname(final String name, final String surname) throws RemoteException {
        System.out.println("Retrieving person " + name);
        for (Person person : persons.values()) {
            if (person.getName().equals(name) && person.getSurname().equals(surname)) {
                return person;
            }
        }
        System.out.println("Person " + name + " " + surname + " not found");
        return null;
    }
}

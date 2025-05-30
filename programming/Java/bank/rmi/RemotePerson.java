package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.rmi.RemoteException;
import java.rmi.server.UnicastRemoteObject;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

public class RemotePerson implements Person {
    private final String name;
    private final String surname;
    private final String passportNumber;
    private final Bank bank;
    private final ConcurrentMap<String, Account> accounts = new ConcurrentHashMap<>();

    /**
     * Creates new person with certain data
     * @param name of person
     * @param surname of person
     * @param passportNumber of person
     * @param bank of person
     * @param port of bank
     * @throws RemoteException If a communication-related exception occurs during RMI call
     */
    public RemotePerson(String name, String surname, String passportNumber, Bank bank, int port) throws RemoteException {
        this.name = name;
        this.surname = surname;
        this.passportNumber = passportNumber;
        this.bank = bank;
        UnicastRemoteObject.exportObject(this, port);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getName() throws RemoteException {
        return name;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getSurname() throws RemoteException {
        return surname;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getPassportNumber() throws RemoteException {
        return passportNumber;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public synchronized Account getAccount(String subId) throws RemoteException {
        String accountId = passportNumber + ":" + subId;
        Account account = accounts.get(subId);
        if (account == null) {
            account = bank.getAccount(accountId);
            if (account != null) {
                accounts.put(subId, account);
            }
        }
        return account;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Account createAccount(String subId) throws RemoteException {
        String accountId = passportNumber + ":" + subId;
        Account account = bank.createAccount(accountId);
        accounts.put(subId, account);
        return account;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public List<String> getAccounts() throws RemoteException {
        return new ArrayList<>(accounts.keySet());
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Bank getBank() throws RemoteException {
        return bank;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void moveAmount(String name, String surname, int amount) throws RemoteException {
        Person person = bank.getPersonByNameAndSurname(name, surname);
        if (person == null) {
            System.out.println("Person " + surname + " " + name + " not found");
            return;
        }
        List<String> friendAccounts = person.getAccounts();
        if (friendAccounts == null) {
            System.out.println("Can't move money to person " + name + " because there are no accounts");
            return;
        }
        if (this.getAccounts() == null || this.getAccounts().isEmpty()) {
            System.out.println("Can't move money from person " + this.getName() + " because there are no accounts");
            return;
        }
        person.getAccount(friendAccounts.getFirst()).addAmount(amount);
        this.getAccount(this.getAccounts().getFirst()).subtractAmount(amount);
        System.out.println("Moved " + amount + " from " + this.getName() + " to " + friendAccounts.getFirst());
    }
}
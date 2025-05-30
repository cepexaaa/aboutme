package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.io.Serializable;
import java.rmi.RemoteException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class LocalPerson implements Person, Serializable {
    private final String name;
    private final String surname;
    private final String passportNumber;
    private final Bank bank;
    private final Map<String, LocalAccount> accounts = new HashMap<>();

    public LocalPerson(Person person) throws RemoteException {
        this.name = person.getName();
        this.surname = person.getSurname();
        this.passportNumber = person.getPassportNumber();
        this.bank = person.getBank();

        for (String subId : person.getAccounts()) {
            Account remoteAccount = person.getAccount(subId);
            accounts.put(subId, new LocalAccount(remoteAccount));
        }
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getName() {
        return name;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getSurname() {
        return surname;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getPassportNumber() {
        return passportNumber;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Account getAccount(String subId) {
        return accounts.get(subId);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Account createAccount(String subId) {
        LocalAccount account = new LocalAccount(passportNumber + ":" + subId, 0);
        accounts.put(subId, account);
        return account;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public List<String> getAccounts() {
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
    public void moveAmount(String name, String surname, int amount) {
        throw new UnsupportedOperationException("LocalPerson cannot modify remote state");
    }
}
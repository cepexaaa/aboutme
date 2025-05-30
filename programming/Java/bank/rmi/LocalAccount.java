package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.io.Serializable;
import java.rmi.RemoteException;

public class LocalAccount implements Account, Serializable {
    private final String id;
    private int amount;

    /**
     * Creates a new local account with the specified ID and balance
     * @param id account identifier
     * @param amount initial balance
     */
    public LocalAccount(String id, int amount) {
        this.id = id;
        this.amount = amount;
    }

    /**
     * Creates a local copy of a remote account
     * @param account remote account to copy data from
     * @throws RemoteException if a remote communication error occurs
     */
    public LocalAccount(Account account) throws RemoteException {
        this.id = account.getId();
        this.amount = account.getAmount();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getId() {
        return id;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public int getAmount() {
        return amount;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void setAmount(int amount) {
        this.amount = amount;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void addAmount(int amount) throws RemoteException {
        setAmount(getAmount() + amount);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void subtractAmount(int amount) throws RemoteException {
        setAmount(getAmount() - amount);
    }
}
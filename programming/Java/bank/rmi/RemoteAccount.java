package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.rmi.RemoteException;

/**
 * {@inheritDoc}
 */
public class RemoteAccount implements Account {
    private final String id;
    private int amount;

    /**
     * Creates new remote account
     * @param id of this new account
     */
    public RemoteAccount(final String id) {
        this.id = id;
        amount = 0;
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
    public synchronized int getAmount() {
        System.out.println("Getting amount of money for account " + id);
        return amount;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public synchronized void setAmount(final int amount) {
        System.out.println("Setting amount of money for account " + id);
        this.amount = amount;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public synchronized void addAmount(int amount) throws RemoteException {
        setAmount(getAmount() + amount);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public synchronized void subtractAmount(int amount) throws RemoteException {
        setAmount(getAmount() - amount);
    }
}

package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.rmi.*;

public interface Account extends Remote {
    /** Returns account identifier. */
    String getId() throws RemoteException;

    /** Returns amount of money in the account. */
    int getAmount() throws RemoteException;

    /** Sets amount of money in the account. */
    void setAmount(int amount) throws RemoteException;

    /** Adding amount of money in the account. */
    void addAmount(int amount) throws RemoteException;

    /** Subtracts amount of money in the account. */
    void subtractAmount(int amount) throws RemoteException;
}

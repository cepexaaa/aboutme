package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.rmi.Remote;
import java.rmi.RemoteException;
import java.util.List;

/**
 * Person in bank app
 * Has main methods to work with accounts and personal data
 */
public interface Person extends Remote {
    /**
     * Returns the person's first name
     * @return first name
     * @throws RemoteException if a remote communication error occurs
     */
    String getName() throws RemoteException;

    /**
     * Returns the person's last name
     * @return last name
     * @throws RemoteException if a remote communication error occurs
     */
    String getSurname() throws RemoteException;

    /**
     * Returns the passport number
     * @return passport number
     * @throws RemoteException if a remote communication error occurs
     */
    String getPassportNumber() throws RemoteException;

    /**
     * Retrieves a bank account by its sub-account identifier
     * @param subId sub-account identifier
     * @return Account object or null if not found
     * @throws RemoteException if a remote communication error occurs
     */
    Account getAccount(String subId) throws RemoteException;

    /**
     * Creates a new bank account with the given sub-account identifier
     * @param subId sub-account identifier
     * @return newly created Account
     * @throws RemoteException if a remote communication error occurs
     */
    Account createAccount(String subId) throws RemoteException;

    /**
     * Returns a list of sub-account identifiers
     * @return list of sub-account IDs
     * @throws RemoteException if a remote communication error occurs
     */
    List<String> getAccounts() throws RemoteException;
    /**
     * Returns a bank of the person
     * @return bank
     * @throws RemoteException if a remote communication error occurs
     */
    Bank getBank() throws RemoteException;

    /**
     * Send amount of money between persons
     * @throws RemoteException if a remote communication error occurs
     */
    void moveAmount(String name, String surname, int amount) throws RemoteException;
}

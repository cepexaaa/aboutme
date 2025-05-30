package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.rmi.Remote;
import java.rmi.RemoteException;

public interface Bank extends Remote {
    /**
     * Creates a new account with specified identifier if it does not already exist.
     * @param id account id
     * @return created or existing account.
     */
    Account createAccount(String id) throws RemoteException;

    /**
     * Returns account by identifier.
     * @param id account id
     * @return account with specified identifier or {@code null} if such account does not exist.
     */
    Account getAccount(String id) throws RemoteException;
    /**
     * Retrieves a remote person reference by their passport number.
     * The returned Person object is a remote stub that will make RMI calls for all operations.
     *
     * @param passportNumber The unique passport identifier of the person
     * @return Remote Person object if found, null if no person exists with this passport
     * @throws RemoteException If a communication-related exception occurs during RMI call
     */
    Person getRemotePerson(String passportNumber) throws RemoteException;
    /**
     * Retrieves a local serializable copy of a person by their passport number.
     * The returned Person object is a local snapshot that won't reflect subsequent changes
     * made by other clients unless refreshed.
     *
     * @param passportNumber The unique passport identifier of the person
     * @return Local Person copy if found, null if no person exists with this passport
     * @throws RemoteException If a communication-related exception occurs during RMI call
     */
    Person getLocalPerson(String passportNumber) throws RemoteException;
    /**
     * Creates a new person record in the bank system.
     *
     * @param name First name of the person
     * @param surname Last name of the person
     * @param passportNumber Unique passport identifier
     * @return The newly created Person object
     * @throws RemoteException If a communication-related exception occurs during RMI call
     * @throws IllegalArgumentException If a person with this passport already exists
     */
    Person createPerson(String name, String surname, String passportNumber) throws RemoteException;
    /**
     * Searches for a person by their first and last name.
     * Note: This may return unexpected results as names are not unique identifiers.
     *
     * @param name First name to search for
     * @param surname Last name to search for
     * @return A Person object if found, null if no matching person exists
     * @throws RemoteException If a communication-related exception occurs during RMI call
     */
    Person getPersonByNameAndSurname(String name, String surname) throws RemoteException;
}

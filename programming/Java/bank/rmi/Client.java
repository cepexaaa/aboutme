package info.kgeorgiy.ja.kubesh.bank.rmi;

import java.net.MalformedURLException;
import java.rmi.Naming;
import java.rmi.NotBoundException;
import java.rmi.RemoteException;

public final class Client {
    /** Utility class. */
    private Client() {}

    /**
     *
     * @param args set data to
     * @throws RemoteException if error would occur
     */
    public static void main(final String... args) throws RemoteException {
        if (args.length < 5) {
            System.out.println("Usage: Client <name> <surname> <passport> <accountId> <amount>");
            return;
        }

        final String name = args[0];
        final String surname = args[1];
        final String passport = args[2];
        final String accountId = args[3];
        final int amount = Integer.parseInt(args[4]);

        final Bank bank;
        try {
            bank = (Bank) Naming.lookup("//localhost/bank");
        } catch (final NotBoundException e) {
            System.out.println("Bank is not bound");
            return;
        } catch (final MalformedURLException e) {
            System.out.println("Bank URL is invalid");
            return;
        }

        Person person = bank.getRemotePerson(passport);
        if (person == null) {
            System.out.println("Creating new person");
            bank.createPerson(name, surname, passport);
        } else {
            System.out.println("Person already exists");
            if (!person.getName().equals(name) || !person.getSurname().equals(surname)) {
                System.out.println("Warning: person data doesn't match");
            }
        }

        Account account = bank.getAccount(accountId);
        if (account == null) {
            System.out.println("Creating account");
            account = bank.createAccount(accountId);
        } else {
            System.out.println("Account already exists");
        }
        System.out.println("Account id: " + account.getId());
        System.out.println("Money: " + account.getAmount());
        System.out.println("Adding money");
        account.setAmount(account.getAmount() + amount);
        System.out.println("Money: " + account.getAmount());
    }
}

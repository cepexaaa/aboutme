import java.util.concurrent.locks.ReentrantLock
import kotlin.concurrent.withLock

/**
 * Bank implementation.
 *
 * @author Кубеш Сергей
 */
class BankImpl(n: Int) : Bank {
    private val accounts: Array<Account> = Array(n) { Account() }

    override val accountsCount: Int
        get() = accounts.size

    override fun amount(id: Int): Long {
        val account = accounts[id]
        account.lock.lock()
            try {
                return accounts[id].amount
            } finally {
                account.lock.unlock()
            }

        }

    override val totalAmount: Long
        get() {
            try {
                for (account in accounts) {
                    account.lock.lock()
                }
                return accounts.sumOf { account ->
                    account.amount
                }
            } finally {
                for (account in accounts.reversed()) {
                    account.lock.unlock()
                }
            }
        }

    override fun deposit(id: Int, amount: Long): Long {
        require(amount > 0) { "Invalid amount: $amount" }
        val account = accounts[id]
//        check(!(amount > Bank.MAX_AMOUNT || account.amount + amount > Bank.MAX_AMOUNT)) { "Overflow" }
        account.lock.lock()
        try {
            check(!(amount > Bank.MAX_AMOUNT || account.amount + amount > Bank.MAX_AMOUNT)) { "Overflow" }
            account.amount += amount
            return account.amount
        } finally {
            account.lock.unlock()
        }
//        account.lock.unlock()
//        return account.amount
    }

    override fun withdraw(id: Int, amount: Long): Long {
        require(amount > 0) { "Invalid amount: $amount" }
        val account = accounts[id]
//        check(account.amount - amount >= 0) { "Underflow" }
//        account.lock.withLock {
//            account.amount -= amount
//            return account.amount
//        }
        account.lock.lock()
        try {
            check(account.amount - amount >= 0) { "Underflow" }
            account.amount -= amount
            return account.amount
        } finally {
            account.lock.unlock()
        }
//        account.amount -= amount
//        account.lock.unlock()
//        return account.amount
    }

    override fun transfer(fromId: Int, toId: Int, amount: Long) {
        require(amount > 0) { "Invalid amount: $amount" }
        require(fromId != toId) { "fromId == toId" }
        val from = accounts[fromId]
        val to = accounts[toId]
//        check(amount <= from.amount) { "Underflow" }
//        check(!(amount > Bank.MAX_AMOUNT || to.amount + amount > Bank.MAX_AMOUNT)) { "Overflow" }
//        to.lock.lock()
//        to.amount += amount
//        to.lock.unlock()
//        from.lock.lock()
//        from.amount -= amount
//        from.lock.unlock()
        val first = if (fromId < toId) from else to
        val second = if (fromId < toId) to else from

        first.lock.lock()
        try {
            second.lock.lock()
            try {
                check(amount <= from.amount) { "Underflow" }
                check(!(amount > Bank.MAX_AMOUNT || to.amount + amount > Bank.MAX_AMOUNT)) { "Overflow" }

                from.amount -= amount
                to.amount += amount
            } finally {
                second.lock.unlock()
            }
        } finally {
            first.lock.unlock()
        }
    }

    /**
     * Private account data structure.
     */
    class Account {
        /**
         * Amount of funds in this account.
         */
        val lock = ReentrantLock()
        var amount: Long = 0
    }
}
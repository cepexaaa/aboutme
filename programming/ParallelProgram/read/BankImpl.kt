import java.util.concurrent.locks.ReentrantReadWriteLock
import kotlin.concurrent.withLock
/**
 * Bank implementation.
 *
 * :TODO: This implementation has to be made thread-safe.
 *
 * @author Кубеш Сергей
 */
class BankImpl(n: Int) : Bank {
    private val accounts: Array<Account> = Array(n) { Account() }

    override val accountsCount: Int
        get() = accounts.size

    /**
     * :TODO: This method has to be made thread-safe.
     */
    override fun amount(id: Int): Long {
        val acc = accounts[id]
        return acc.lock.readLock().withLock {
            acc.amount
        }
    }

    /**
     * :TODO: This method has to be made thread-safe.
     */
    override val totalAmount: Long
        get() {
            var sum: Long = 0
//            for (i in 0 until accounts.size) {
//                sum+=amount(i)
//            }
            for (i in accounts.indices) {
                accounts[i].lock.readLock().lock()
            }
            try {
                for (account in accounts) {
                    sum += account.amount
                }
            } finally {
                for (i in accounts.indices.reversed()) {
                    accounts[i].lock.readLock().unlock()
                }
            }
//            for (account in accounts) {
//                account.lock.readLock().withLock {
//                    sum += account.amount
//                }
//            }
            return sum
//            accounts.sumOf { account ->
//                account.amount
//            }
        }

    /**
     * :TODO: This method has to be made thread-safe.
     */
    override fun deposit(id: Int, amount: Long): Long {
        require(amount > 0) { "Invalid amount: $amount" }
        val account = accounts[id]
        return account.lock.writeLock().withLock {
            check(amount <= Bank.MAX_AMOUNT && account.amount + amount <= Bank.MAX_AMOUNT) { "Overflow" }
            account.amount += amount
            account.amount
        }
    }

    /**
     * :TODO: This method has to be made thread-safe.
     */
    override fun withdraw(id: Int, amount: Long): Long {
        require(amount > 0) { "Invalid amount: $amount" }
        val account = accounts[id]
        return account.lock.writeLock().withLock {
            check(account.amount - amount >= 0) { "Underflow" }
            account.amount -= amount
            account.amount
        }
    }

    /**
     * :TODO: This method has to be made thread-safe.
     */
    override fun transfer(fromId: Int, toId: Int, amount: Long) {
        require(amount > 0) { "Invalid amount: $amount" }
        require(fromId != toId) { "fromId == toId" }
//        val from = accounts[fromId]
//        val to = accounts[toId]

        val (firstId, secondId) = if (fromId < toId) fromId to toId else toId to fromId
        val first = accounts[firstId] //if (fromId < toId) from else to
        val second = accounts[secondId] //if (fromId < toId) to else from
        first.lock.writeLock().withLock {
            second.lock.writeLock().withLock {
                val from = accounts[fromId]
                val to = accounts[toId]

                check(amount <= from.amount) { "Underflow" }
                check(to.amount + amount <= Bank.MAX_AMOUNT) { "Overflow" }

                from.amount -= amount
                to.amount += amount
            }
        }
//        val l1 = first.lock.writeLock()
//        try {
//            val l2 = second.lock.writeLock()
//            try {
//                check(amount <= from.amount) { "Underflow" }
//                check(amount <= Bank.MAX_AMOUNT && to.amount + amount <= Bank.MAX_AMOUNT) { "Overflow" }
//                from.amount -= amount
//                to.amount += amount
//            } finally {
//                l2.unlock()
//            }
//        } finally {
//            l1.unlock()
//        }
    }

    /**
     * :TODO: This method has to be made thread-safe.
     */
    override fun consolidate(fromIds: List<Int>, toId: Int) {
        require(fromIds.isNotEmpty()) { "empty fromIds" }
        require(fromIds.distinct() == fromIds) { "duplicates in fromIds" }
        require(toId !in fromIds) { "toId in fromIds" }
        val sortedInd = (fromIds + toId).sorted()
//        val accountsToLock = sortedInd.map { accounts[it] }
//        val to = accounts[toId]
//        for (account in accountsToLock) {
//            account.lock.writeLock().lock()
//        }
        for (id in sortedInd) {
            accounts[id].lock.writeLock().lock()
        }
//        val locks: Array<ReentrantReadWriteLock.WriteLock> = Array(sortedInd.size) { ReentrantReadWriteLock.WriteLock() }
//        for (i in 0 until accountsToLock.size) {
//            locks[i] = accountsToLock[i].lock.writeLock()
//        }
        try {
            val amount = fromIds.sumOf { accounts[it].amount }
            check(accounts[toId].amount + amount <= Bank.MAX_AMOUNT) { "Overflow" }
            for (from in fromIds) { accounts[from].amount = 0}
            accounts[toId].amount += amount
        } finally {
            for (id in sortedInd.reversed()) {
                accounts[id].lock.writeLock().unlock()
            }
//            for (account in accountsToLock.reversed()) {
//                account.lock.writeLock().unlock()
//            }
        }
//        for (l in locks) {
//            l.unlock()
//        }
    }

    /**
     * Private account data structure.
     */
    class Account {
        /**
         * Amount of funds in this account.
         */
        val lock = ReentrantReadWriteLock()
        var amount: Long = 0
    }
}
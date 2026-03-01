import kotlin.concurrent.atomics.AtomicArray
import kotlin.concurrent.atomics.ExperimentalAtomicApi

/**
 * Bank implementation.
 * This class is thread-safe and lock-free using operation objects.
 *
 * This implementation is based on "A Practical Multi-Word Compare-and-Swap Operation" by T. L. Harris et al.
 * It uses a simplified and faster version of DCSS operation that relies for its correctness on the fact that
 * Account instances in [accounts] array never suffer from ABA problem.
 * See also "Practical lock-freedom" by Keir Fraser. See [acquire] method.
 *
 * :TODO: This implementation has to be completed, so that it is thread-safe and lock-free.
 */
@OptIn(ExperimentalAtomicApi::class)
class BankImpl(override val accountsCount: Int) : Bank {
    /**
     * An array of accounts.
     * Account instances here are never reused (there is no ABA).
     */
    private val accounts = AtomicArray<Account>(accountsCount) { Account(0) }

    private fun account(id: Int) = accounts.loadAt(id)

    override fun amount(id: Int): Long {
        while (true) {
            val account = account(id)
            /*
             * If there is a pending operation on this account, then help to complete it first using
             * its invokeOperation method. If the result is false then there is no pending operation,
             * thus the account amount can be safely returned.
             */
            if (!account.invokeOperation()) return account.amount
        }
    }

    override val totalAmount: Long
        get() {
            /**
             * This operation requires atomic read of all accounts, thus it creates an operation descriptor.
             * Operation's invokeOperation method acquires all accounts, computes the total amount, and releases
             * all accounts. This method returns the result.
             */
            val op = TotalAmountOp()
            op.invokeOperation()
            return op.sum
        }

    override fun deposit(id: Int, amount: Long): Long { // First, validate method per-conditions
        require(amount > 0) { "Invalid amount: $amount" }
        check(amount <= MAX_AMOUNT) { "Overflow" }
        /*
         * This operation depends only on a single account, thus it can be directly
         * performed using a regular lock-free compareAndSet loop.
         */
        while (true) {
            val account = account(id)
            /*
             * If there is a pending operation on this account, then help to complete it first using
             * its invokeOperation method. If the result is false then there is no pending operation,
             * thus the account can be safely updated.
             */
            if (account.invokeOperation()) continue
            check(account.amount + amount <= MAX_AMOUNT) { "Overflow" }
            val updated = Account(account.amount + amount)
            if (accounts.compareAndSetAt(id, account, updated)) return updated.amount
        }
    }

    override fun withdraw(id: Int, amount: Long): Long {
        // todo: write withdraw operation using deposit as an example
        /*
         * Basically, implementation of this method must perform the logic of the following code "atomically":
         */
        require(amount > 0) { "Invalid amount: $amount" }
        while (true) {
            val account = account(id)
            if (account.invokeOperation()) continue
            check(account.amount - amount >= 0) { "Underflow" }
            val updated = Account(account.amount - amount)
            if(accounts.compareAndSetAt(id, account, updated)) return updated.amount
        }
//        require(amount > 0) { "Invalid amount: $amount" }
//        val account = account(id)
//        check(account.amount - amount >= 0) { "Underflow" }
//        val updated = Account(account.amount - amount)
//        accounts.storeAt(id, updated)
//        return updated.amount
    }

    override fun transfer(fromId: Int, toId: Int, amount: Long) {
        // First, validate method per-conditions
        require(amount > 0) { "Invalid amount: $amount" }
        require(fromId != toId) { "fromId == toId" }
        check(amount <= MAX_AMOUNT) { "Underflow/overflow" }
        /**
         * This operation requires atomic read of two accounts, thus it creates an operation descriptor.
         * Operation's invokeOperation method acquires both accounts, computes the result of operation
         * (if a form of error message), and releases both accounts. This method throws the exception with
         * the corresponding message if needed.
         */
        val op = TransferOp(fromId, toId, amount)
        op.invokeOperation()
        op.errorMessage?.let { error(it) }
    }

    /**
     * This is an implementation of a restricted form of Harris DCSS operation:
     * It atomically checks that op.completed is false and replaces accounts[id] with AcquiredAccount instance
     * that hold a reference to the op.
     * This method returns null if op.completed is true.
     */
    private fun acquire(id: Int, op: Op): AcquiredAccount? {
        // todo: write the implementation of this method with the following logic:
        /*
         * This method must loop trying to replace accounts[id] with an instance of
         *     new AcquiredAccount(<old-amount>, op) until that successfully happens and return the
         *     instance of AcquiredAccount in this case.
         *
         * If current account is already "Acquired" by another operation, then this method must help that
         * other operation by invoking "invokeOperation" and continue trying.
         *
         * Because accounts[id] does not have an ABA problem, there is no need to implement full-blown
         * DCSS operation with descriptors for DCSS operation as explained in Harris CASN work. A simple
         * lock-free compareAndSet loop suffices here.
         *
         * **Here are the important details:**
         *  - op.completed must be checked **after** the accounts[id] is read.
         *  - check if the account is **already acquire** for this op and return the current AcquiredAccount if so.
         *
         * Basically, implementation of this method must perform the logic of the following code "atomically":
         */
        while (true) {
//            val account = account(id)
//            if (a is AcquiredAccount) return a
            val account = accounts.loadAt(id)
            if (op.completed) return null

            if (account is AcquiredAccount && account.op === op) {
                return account
            }
//            if (account is AcquiredAccount) {
//                account.invokeOperation()
//                if (!account.op.completed) {
//                    continue
//                }
//            }
            if (!account.invokeOperation()) {
                val acquiredAccount = AcquiredAccount(account.amount, op)
                if (accounts.compareAndSetAt(id, account, acquiredAccount)) {
                    return acquiredAccount
                }
            }
        }
//        if (op.completed) return null
//        val account = account(id)
//        val acquiredAccount = AcquiredAccount(account.amount, op)
//        accounts.storeAt(id, acquiredAccount)
//        return acquiredAccount
    }

    /**
     * Releases an account that was previously acquired by [acquire].
     * This method does nothing if the account at id is not currently acquired.
     */
    private fun release(id: Int, op: Op) {
        assert(op.completed) // must be called only on operations that were already completed
        val account = account(id)
        if (account is AcquiredAccount && account.op === op) {
            // release performs update at most once while the account is still acquired
            val updated = Account(account.newAmount)
            accounts.compareAndSetAt(id, account, updated)
        }
    }

    /**
     * Immutable account data structure.
     * @param amount Amount of funds in this account.
     */
    private open class Account(val amount: Long) {
        /**
         * Invokes operation that is pending on this account.
         * This implementation returns false (no pending operation), other implementations return true.
         */
        open fun invokeOperation(): Boolean = false
    }

    /**
     * Account that was acquired as a part of in-progress operation that spans multiple accounts.
     * @see acquire
     */
    private class AcquiredAccount(
        var newAmount: Long, // New amount of funds in this account when op completes.
        val op: Op
    ) : Account(newAmount) {
        override fun invokeOperation(): Boolean {
            op.invokeOperation()
            return true
        }
    }

    /**
     * Abstract operation that acts on multiple accounts.
     */
    private abstract class Op {
        /**
         * True when operation has completed.
         */
        @Volatile
        var completed = false

        abstract fun invokeOperation()
    }

    /**
     * Descriptor for [totalAmount] operation.
     */
    private inner class TotalAmountOp : Op() {
        /**
         * The result of getTotalAmount operation is stored here before setting
         * [completed] to true.
         */
        var sum = 0L

        override fun invokeOperation() {
            var sum = 0L
            var acquired = 0
            while (acquired < accountsCount) {
                val account = acquire(acquired, this) ?: break
                sum += account.newAmount
                acquired++
            }
            if (acquired == accountsCount) {
                /*
                 * If i == n, then all acquired accounts were not null and full sum was calculated.
                 * this.sum = sum assignment below has a benign data race. Multiple threads might to this assignment
                 * concurrently, however, they are all guaranteed to be assigning the same value.
                 */
                this.sum = sum
                completed = true // volatile write to completed field _after_ the sum was written
            }
            /*
             * To ensure lock-freedom, we must release all accounts even if this particular helper operation
             * had failed to acquire all of them before somebody else had completed the operations.
             * By releasing all accounts for completed operation we ensure progress of other operations.
             */
            for (i in 0 until accountsCount) {
                release(i, this)
            }
        }
    }

    /**
     * Descriptor for [transfer] operation.
     */
    private inner class TransferOp(val fromId: Int, val toId: Int, val amount: Long) : Op() {
        var errorMessage: String? = null
//        var checked = false

        override fun invokeOperation() {
            // todo: write implementation for this method, use TotalAmountOp as an example
            /*
             * In the implementation of this operation only two accounts (with fromId and toId) needs
             * to be acquired. Unlike TotalAmountOp, this operation has its own result in errorMessage string,
             * and it must also update AcquiredAccount.newAmount fields before setting completed to true
             * and invoking release on those acquired accounts.
             *
             * Basically, implementation of this method must perform the logic of the following code "atomically":
             */

//            if (completed) {
//                release(fromId, this)
//                release(toId, this)
//                return
//            }

            acquire(min(fromId, toId), this)
            acquire(max(fromId, toId), this)


            val fromAccount = acquire(fromId, this)
            val toAccount = acquire(toId, this)

            if (fromAccount != null && toAccount != null) {
                val fromNewAmount = fromAccount.amount - amount
                val toNewAmount = toAccount.amount + amount

//                if (!completed/* && !checked*/) {
                    when {
                        fromNewAmount < 0 -> errorMessage = "Underflow"
                        toNewAmount > MAX_AMOUNT -> errorMessage = "Overflow"
                        else -> {
//                            checked = true
                            fromAccount.newAmount = fromNewAmount
                            toAccount.newAmount = toNewAmount

//                            release(fromId, this)
//                            release(toId, this)
                        }
                    }
                }

//            if (!completed) {
            this.completed = true
            release(fromId, this)
            release(toId, this)

//            }
//            val from = account(fromId)
//            val to = account(toId)
//            when {
//                amount > from.amount -> errorMessage = "Underflow"
//                to.amount + amount > MAX_AMOUNT -> errorMessage = "Overflow"
//                else -> {
//                    accounts.storeAt(fromId, Account(from.amount - amount))
//                    accounts.storeAt(toId, Account(to.amount + amount))
//                }
//            }
        }
    }

    fun min(a :Int, b: Int) = if (a <= b) a else b
    fun max(a :Int, b: Int) = if (a >= b) a else b
}
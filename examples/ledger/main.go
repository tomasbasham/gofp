// Ledger is ledger system implemented using a [state.State] monad. The ledger
// supports deposit, withdraw, and transaction operations, while maintaining a
// running total balance.
package main

import (
	"crypto/rand"
	"fmt"

	"github.com/tomasbasham/gofp"
	"github.com/tomasbasham/gofp/state"
)

// EntryID is a unique identifier for a ledger entry.
type EntryID string

// Ledger represents a simple ledger system mapping entry IDs to ledger entries.
type Ledger map[EntryID]LedgerEntry

// LedgerEntry represents a single entry in a ledger.
type LedgerEntry struct {
	Amount      int
	Description string

	// The balance after this entry has been applied. It is used to keep a running
	// total of the balance of the ledger for performance.
	balanceAfter int
}

// LedgerState is a state monad that operates on a ledger.
type LedgerState = state.State[Ledger, gofp.Result[int]]

// Operation represents a transaction operation.
type Operation int8

const (
	Deposit Operation = iota
	Withdraw
	Transaction
)

func (o Operation) String() string {
	switch o {
	case Deposit:
		return "deposit"
	case Withdraw:
		return "withdraw"
	case Transaction:
		return "transaction"
	default:
		return "unknown"
	}
}

// TransactionError is an error that occurs during a transaction.
type TransactionError struct {
	code    string
	message string
}

func (e *TransactionError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

// NewInsufficientFundsError creates a new TransactionError for insufficient
// funds.
func NewInsufficientFundsError(balance, amount int) *TransactionError {
	return &TransactionError{
		code:    "INSUFFICIENT_FUNDS",
		message: fmt.Sprintf("insufficient funds for withdrawal: balance %d, requested %d", balance, amount),
	}
}

// NewInvalidAmountError creates a new TransactionError for an invalid amount.
func NewInvalidAmountError(amount int) *TransactionError {
	return &TransactionError{
		code:    "INVALID_AMOUNT",
		message: fmt.Sprintf("invalid amount: %d", amount),
	}
}

func main() {
	s := newLedger()

	// Deposit and withdraw some money.
	s, d1 := deposit(s, 75)
	s, w1 := withdraw(s, 50)

	// Create a new transaction and deposit and withdraw some money.
	tx, t1 := newTransaction()
	tx, d2 := deposit(tx, 100)
	tx, w2 := withdraw(tx, 50)

	// Commit the transaction to the base ledger. With this design it is possible
	// to create many nested transactions as the basis for this implementation is
	// to apply state on top of other states.
	s = commit(s, tx)

	// Get the balance of the ledger. This effectively applies all of the state
	// operations to the ledger and returns the final balance, and an updated
	// ledger.
	balance, ledger := s.Run(Ledger{})
	if balance.IsErr() {
		fmt.Printf("Error: %v\n", balance.UnwrapErr())
		return
	}

	// Create and commit another transaction. This demonstrates that even after
	// applying previous operations to the ledger, it is still possible to create
	// new transactions and apply them to the ledger.
	tx2, t2 := newTransaction()
	tx2, d3 := deposit(tx2, 60)
	tx2, w3 := withdraw(tx2, 10)

	// Commit the second transaction.
	s = commit(s, tx2)

	balance, _ = s.Run(ledger)
	if balance.IsErr() {
		fmt.Printf("Error: %v\n", balance.UnwrapErr())
		return
	}

	// Print the final balance and some balances at specific points in the ledger.
	// Since the ledger is not more than a map of entries, it is possible to query
	// the ledger at any point in time.
	fmt.Printf("Final balance: %d\n", balance.Unwrap())
	printBalances(ledger, d1, w1, t1, d2, w2, t2, d3, w3)
}

func newLedger() LedgerState {
	return state.Pure[Ledger](gofp.Ok(0))
}

func deposit(s LedgerState, amount int) (LedgerState, EntryID) {
	if amount <= 0 {
		return state.Pure[Ledger](gofp.Err[int](NewInvalidAmountError(amount))), ""
	}
	return update(s, amount, Deposit, func(balance int) gofp.Result[int] {
		return gofp.Ok(balance + amount)
	})
}

func withdraw(s LedgerState, amount int) (LedgerState, EntryID) {
	if amount <= 0 {
		return state.Pure[Ledger](gofp.Err[int](NewInvalidAmountError(amount))), ""
	}
	return update(s, -amount, Withdraw, func(balance int) gofp.Result[int] {
		if balance < amount {
			return gofp.Err[int](NewInsufficientFundsError(balance, amount))
		}
		return gofp.Ok(balance - amount)
	})
}

func newTransaction() (LedgerState, EntryID) {
	return update(newLedger(), 0, Transaction, func(balance int) gofp.Result[int] {
		return gofp.Ok(balance)
	})
}

func commit(s1 LedgerState, s2 LedgerState) LedgerState {
	return state.Zip(s1, s2, func(l1, l2 gofp.Result[int]) gofp.Result[int] {
		return l1.AndThen(func(b1 int) gofp.Result[int] {
			return l2.AndThen(func(b2 int) gofp.Result[int] {
				return gofp.Ok(b1 + b2)
			})
		})
	})
}

func update(s LedgerState, amount int, op Operation, fn func(int) gofp.Result[int]) (LedgerState, EntryID) {
	id := newEntryID(op)
	return s.FlatMap(func(res gofp.Result[int]) LedgerState {
		if res.IsErr() {
			return state.Pure[Ledger](res)
		}

		newBalance := fn(res.Unwrap())
		if newBalance.IsErr() {
			return state.Pure[Ledger](newBalance)
		}

		balanceAfter := newBalance.Unwrap()
		return state.FlatMap(state.Get[Ledger](), func(l Ledger) LedgerState {
			desc := fmt.Sprintf("%s of %d", op, amount)
			if op == Transaction {
				desc = "open transaction"
			}
			l[id] = LedgerEntry{
				Amount:       amount,
				balanceAfter: balanceAfter,
				Description:  desc,
			}
			return state.FlatMap(state.Put(l), func(gofp.Unit) LedgerState {
				return state.Pure[Ledger](gofp.Ok(balanceAfter))
			})
		})
	}), id
}

func newEntryID(op Operation) EntryID {
	if op == Transaction {
		return EntryID(newID("tx"))
	}
	return EntryID(newID("ent"))
}

func newID(prefix string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s-%x", prefix, b[0:8])
}

func balanceAt(l Ledger, id EntryID) gofp.Option[int] {
	if entry, ok := l[id]; ok {
		return gofp.Some(entry.balanceAfter)
	}
	return gofp.None[int]()
}

func printBalances(l Ledger, ids ...EntryID) {
	for _, id := range ids {
		b := balanceAt(l, id)
		if b.IsSome() {
			desc := descriptionAt(l, id).UnwrapOr("No description")
			fmt.Printf("Balance at %s: %d (%s)\n", id, b.Unwrap(), desc)
		}
	}
}

func descriptionAt(l Ledger, id EntryID) gofp.Option[string] {
	if entry, ok := l[id]; ok {
		return gofp.Some(entry.Description)
	}
	return gofp.None[string]()
}

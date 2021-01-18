package account

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

// constanst for trx types of debit/credit
const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

// TransactionType is credit or debit
type TransactionType string

// IsValid check if trx type is in the enum
func (t TransactionType) IsValid() bool {
	return t == Credit || t == Debit
}

// IsCredit check if is credit
func (t TransactionType) IsCredit() bool {
	return t == Credit
}

// IsCredit check if is debit
func (t TransactionType) IsDebit() bool {
	return t == Debit
}

// Transaction type
type Transaction struct {
	id            big.Int
	kind          TransactionType
	amount        decimal.Decimal
	effectiveDate time.Time
}

// MarshalJSON converts to json format
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID            string `json:"id"`
		Kind          string `json:"type"`
		Amount        string `json:"amount"`
		EffectiveDate string `json:"effective_date"`
	}{
		ID:            t.id.String(),
		Kind:          string(t.kind),
		Amount:        t.amount.StringFixedBank(2),
		EffectiveDate: t.effectiveDate.Format(time.ANSIC),
	})
}

// Transactions array type
type Transactions []Transaction

// History list all the transactions
var History = make(Transactions, 0)

// CreateTx do a transaction and returns the transaction
func CreateTx(kind TransactionType, amount decimal.Decimal) (tx Transaction, err error) {

	Status.insertMu.Lock()
	defer Status.insertMu.Unlock()

	amount = amount.Abs()

	newBalance := decimal.New(0, 10)
	switch kind {
	case Credit:
		newBalance = Status.balance.Add(amount)
	case Debit:
		newBalance = Status.balance.Sub(amount)
	default:
		return tx, fmt.Errorf("tx type invalid")
	}

	if newBalance.IsNegative() {
		return tx, fmt.Errorf("balance is negative")
	}

	tx = Transaction{
		id:     *Status.count.Add(&Status.count, big.NewInt(1)),
		amount: amount,
		kind:   kind,
	}
	History = append(History, tx)

	Status.balance = newBalance
	Status.count = tx.id

	return tx, nil
}

// Find for a transaction in the list
func (list Transactions) Find(id big.Int) (tx *Transaction) {
	Status.insertMu.RLock()
	defer Status.insertMu.RUnlock()

	for _, tx := range list {
		if id.String() == tx.id.String() {
			return &tx
		}
	}

	return nil
}

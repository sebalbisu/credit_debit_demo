package account

import (
	"math/big"
	"sync"

	"github.com/shopspring/decimal"
)

// Status of the bank account
var Status = BankStatus{
	count:   *big.NewInt(0),
	balance: decimal.NewFromFloat(0),
}

// BankStatus account of the user
type BankStatus struct {
	count    big.Int
	balance  decimal.Decimal
	insertMu sync.RWMutex
}

// Count of trx user has made, and last insert id
func (s *BankStatus) Count() big.Int {
	s.insertMu.RLock()
	defer s.insertMu.RUnlock()

	return s.count
}

// Balance account of the user
func (s *BankStatus) Balance() decimal.Decimal {
	s.insertMu.RLock()
	defer s.insertMu.RUnlock()

	return s.balance
}

// BalanceFormated  of the account of the user
func (s *BankStatus) BalanceFormated() string {
	s.insertMu.RLock()
	defer s.insertMu.RUnlock()

	return s.balance.StringFixedBank(2)
}

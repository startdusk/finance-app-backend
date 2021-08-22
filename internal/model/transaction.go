package model

import (
	"errors"
	"time"
)

// TransactionID is identifier of Transaction
type TransactionID string

// NilTransactionID is empty identifier for Transaction
var NilTransactionID TransactionID

// TransactionType is string representation of Transaction type
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID         TransactionID `json:"id" db:"transaction_id"`
	UserID     *UserID       `json:"userID" db:"user_id"`
	AccountID  *AccountID    `json:"accountID" db:"account_id"`
	CategoryID *CategoryID   `json:"categoryID" db:"category_id"`

	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`

	Date   *time.Time       `json:"date" db:"date"`
	Type   *TransactionType `json:"type" db:"type"`
	Amount *int64           `json:"amount" db:"amount"`
	Notes  *string          `json:"notes" db:"notes"`
}

func (t *Transaction) Verify() error {
	if t.UserID == nil || len(*t.UserID) == 0 {
		return errors.New("userID is required")
	}

	if t.AccountID == nil || len(*t.AccountID) == 0 {
		return errors.New("accountID is required")
	}

	if t.CategoryID == nil || len(*t.CategoryID) == 0 {
		return errors.New("categoryID is required")
	}

	if t.Date == nil {
		return errors.New("date is required")
	}

	if t.Type == nil || len(*t.Type) == 0 {
		return errors.New("type is required")
	}

	if t.Amount == nil {
		return errors.New("amount is required")
	}

	return nil
}

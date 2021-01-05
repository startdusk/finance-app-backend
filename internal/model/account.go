package model

import (
	"errors"
	"time"
)

// AccountID is an identifier for account
type AccountID string

// NilAccountID is empty AccountID
var NilAccountID AccountID

// AccountType is type of account
type AccountType string

const (
	Cash   AccountType = "cash"
	Credit AccountType = "credit"
)

// Account is structure for account
type Account struct {
	ID           AccountID    `json:"id,omitempty" db:"account_id"`
	UserID       *UserID      `json:"userID,omitempty" db:"user_id"`
	Name         *string      `json:"name,omitempty" db:"account_name"`
	Type         *AccountType `json:"type,omitempty" db:"account_type"`
	StartBalance *int64       `json:"startBalance,omitempty" db:"start_balance"`
	Currency     *string      `json:"currency,omitempty" db:"currency"`
	CreatedAt    *time.Time   `json:"-" db:"created_at"`
	DeletedAt    *time.Time   `json:"-" db:"deleted_at"`
}

func (a *Account) Verify() error {
	if a.UserID == nil || len(*a.UserID) == 0 {
		return errors.New("UserID is required")
	}

	if a.Name == nil || len(*a.Name) == 0 {
		return errors.New("Name is required")
	}

	if a.Type == nil || len(*a.Type) == 0 {
		return errors.New("Type is required")
	}

	if a.StartBalance == nil {
		return errors.New("StartBalance is required")
	}

	if a.Currency == nil || len(*a.Currency) == 0 {
		return errors.New("Currency is required")
	}

	return nil
}

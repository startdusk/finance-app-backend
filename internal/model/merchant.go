package model

import (
	"errors"
	"time"
)

// MerchantID is identifier of Merchant
type MerchantID string

// NilMerchantID is empty identifier of Merchant
var NilMerchantID MerchantID

type Merchant struct {
	ID        MerchantID `json:"id,omitempty" db:"merchant_id"`
	UserID    *UserID    `json:"userID,omitempty" db:"user_id"`
	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	Name      *string    `json:"name,omitempty" db:"name"`
}

func (a *Merchant) Verify() error {
	if a.UserID == nil || len(*a.UserID) == 0 {
		return errors.New("userID is required")
	}

	if a.Name == nil || len(*a.Name) == 0 {
		return errors.New("name is required")
	}

	return nil
}

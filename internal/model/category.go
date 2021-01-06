package model

import (
	"errors"
	"time"
)

// CategoryID is identifier of Category
type CategoryID string

// NilCategoryID is empty identifier of Category
var NilCategoryID CategoryID

type Category struct {
	ID        CategoryID `json:"id,omitempty" db:"category_id"`
	ParentID  CategoryID `json:"parentID,omitempty" db:"parent_id"`
	UserID    *UserID    `json:"userID,omitempty" db:"user_id"`
	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	Name      *string    `json:"name,omitempty" db:"name"`
}

func (a *Category) Verify() error {
	if a.UserID == nil || len(*a.UserID) == 0 {
		return errors.New("userID is required")
	}

	if a.Name == nil || len(*a.Name) == 0 {
		return errors.New("name is required")
	}

	return nil
}

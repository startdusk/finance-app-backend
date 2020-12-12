package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserID is identifier for User
type UserID string

// User is structure represent ther object
type User struct {
	ID           UserID     `json:"id,omitempty" db:"user_id"`
	Email        *string    `json:"email" db:"email"`
	PasswordHash *[]byte    `json:"-" db:"password_hash"`
	CreatedAt    *time.Time `json:"-" db:"created_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

// Verify all required fields before create or update
func (u *User) Verify() error {
	if u.Email == nil || (u.Email != nil && len(*u.Email) == 0) {
		return errors.New("Email is required")
	}

	return nil
}

// SetPassword updates user's password
func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.PasswordHash = &hash
	return nil
}

// CheckPassword verifies user's password
func (u *User) CheckPassword(password string) error {
	if u.PasswordHash != nil && len(*u.PasswordHash) == 0 {
		return errors.New("password not set")
	}
	return bcrypt.CompareHashAndPassword(*u.PasswordHash, []byte(password))
}

// HashPassword hashes a user's raw password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

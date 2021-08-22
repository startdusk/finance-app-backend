package model

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserID is identifier for User
type UserID string

// NilUserID is an empty UserID
var NilUserID UserID

// User is structure represent ther object
type User struct {
	ID           UserID     `json:"id,omitempty" db:"user_id"`
	Email        *string    `json:"email" db:"email"`
	PasswordHash *[]byte    `json:"-" db:"password_hash"`
	CreatedAt    *time.Time `json:"-" db:"created_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

// isEmail match email format
var isEmail = regexp.MustCompile(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`)

// Verify all required fields before create or update
func (u *User) Verify() error {
	if u.Email == nil || (u.Email != nil && len(*u.Email) == 0) {
		return errors.New("email is required")
	}
	if !isEmail.MatchString(*u.Email) {
		return errors.New("email invalid")
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

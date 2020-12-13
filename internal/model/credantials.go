package model

import (
	"fmt"
)

// Credantials used in login API
type Credantials struct {

	// Username/Password login:
	Email    string `json:"email"`
	Password string `json:"password"`

	// In future we will have google and facebook login as well
}

// Principal is an authenticated entity
type Principal struct {
	UserID UserID `json:"userID,omitempty"`
}

func (p Principal) String() string {
	if p.UserID != "" {
		return fmt.Sprintf("UserID[%s]", p.UserID)
	}
	return "(none)"
}

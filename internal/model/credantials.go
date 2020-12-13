package model

// Credantials used in login API
type Credantials struct {

	// Username/Password login:
	Email    string `json:"email"`
	Password string `json:"password"`

	// In future we will have google and facebook login as well
}

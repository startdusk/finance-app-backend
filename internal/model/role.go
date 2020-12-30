package model

// Role is a function a user can serve
type Role string

const (
	// RoleAdmin is an administrator of App. Root
	RoleAdmin Role = "admin"
)

// We need some structure which will be represent Role
// It can be list of Role(string) but whta if we want to add some aditional info in future...
// So let's make it structure
type UserRole struct {
	Role Role `json:"role" db:"role"`
}

package auth

import (
	"github.com/startdusk/finance-app-backend/internal/model"
)

// We will have 3 permission type for now.
type PermissionType string

const (
	// User has 'admin' role
	Admin PermissionType = "admin"
	// User is loged in (we have user id in principal)
	Member PermissionType = "member"
	// User is loged in and user id passed to API is the same
	MemberIsTarget PermissionType = "memberIsTarget"
	// Any one can access
	Any PermissionType = "anonym"
)

// We will create function for each type

// Admin
var adminOnly = func(roles []*model.UserRole) bool {
	for _, role := range roles {
		switch role.Role {
		case model.RoleAdmin:
			return true
		}
	}
	return false
}

// Loged in user(用户已经登陆了，token里面携带了userID)
var member = func(principal model.Principal) bool {
	return principal.UserID != ""
}

// Loged in user - Target user(path带userID和token里面带的userID必须是一致，防止越权操作)
var memberIsTarget = func(userID model.UserID, principal model.Principal) bool {
	if userID == "" || principal.UserID == "" {
		return false
	}

	if userID != principal.UserID {
		return false
	}

	return true
}

var any = func() bool {
	return true
}

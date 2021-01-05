package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bluele/gcache"
	"github.com/gorilla/mux"

	"github.com/startdusk/finance-app-backend/internal/api/utils"
	"github.com/startdusk/finance-app-backend/internal/database"
	"github.com/startdusk/finance-app-backend/internal/model"
)

type Permissions interface {
	Wrap(next http.HandlerFunc, permissionTypes ...PermissionType) http.HandlerFunc
	Check(r *http.Request, permissionTypes ...PermissionType) bool
}

type permissions struct {
	DB    database.Database
	cache gcache.Cache
}

func NewPermissions(db database.Database) Permissions {
	p := &permissions{
		DB: db,
	}

	p.cache = gcache.New(20).
		LRU().
		LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
			userID, ok := key.(model.UserID)
			if !ok {
				return nil, nil, fmt.Errorf("unknow key type: %v", key)
			}
			roles, err := p.DB.GetRolesByUser(context.Background(), userID)
			if err != nil {
				return nil, nil, err
			}
			expire := 1 * time.Minute
			return roles, &expire, nil
		}).
		Build()

	return p
}

// get user's roles from cache (if we wont have roles in cache it will get it from database)
func (p *permissions) getRoles(userID model.UserID) ([]*model.UserRole, error) {
	r, err := p.cache.Get(userID)
	if err != nil {
		return nil, err
	}

	roles, ok := r.([]*model.UserRole)
	if !ok {
		return nil, fmt.Errorf("cannot get roles: %v", roles)
	}
	return roles, nil
}

func (p *permissions) withRoles(principal model.Principal, roleFunc func([]*model.UserRole) bool) (bool, error) {
	if principal.UserID == model.NilUserID {
		return false, nil
	}

	// Load roles
	roles, err := p.getRoles(principal.UserID)
	if err != nil {
		return false, err
	}

	return roleFunc(roles), nil
}

// we need see if we have principal on Request in this point...
func (p *permissions) Wrap(next http.HandlerFunc, permissionTypes ...PermissionType) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowed := p.Check(r, permissionTypes...); !allowed {
			utils.WriteError(w, http.StatusUnauthorized, "permission denied", nil)
			return
		}
		// also we need to check if we can get userid from path in case we want to check if authorized user is user who trying access api
		// for example user can get only himself
		// principal := GetPrincipal(r)
		// fmt.Println("Test from Permissions principal:", principal)

		// one more test...
		// I want to check if we get any errors if we wrap API where we don't pass userID
		// vars := mux.Vars(r)
		// userID := model.UserID(vars["userID"])
		// fmt.Println("Test from Permissions userID:", userID)

		next.ServeHTTP(w, r)
	})
}

// The idea is to return TRUE if one of permission types matches.
// for example If permission type is Admin and MemberIsTarget
// Admin can edit any user so if user has Admin role we don't care, admin don't match MemberIsTarget permission
func (p *permissions) Check(r *http.Request, permissionTypes ...PermissionType) bool {
	principal := GetPrincipal(r)
	for _, permissionType := range permissionTypes {
		switch permissionType {
		case Admin:
			if allowed, _ := p.withRoles(principal, adminOnly); allowed {
				return true
			}
		case Member:
			if allowed := member(principal); allowed {
				return true
			}
		case MemberIsTarget:
			targetUserID := model.UserID(mux.Vars(r)["userID"])
			if allowed := memberIsTarget(targetUserID, principal); allowed {
				return true
			}
		case Any:
			if allowed := any(); allowed {
				return true
			}
		}
	}
	return false
}

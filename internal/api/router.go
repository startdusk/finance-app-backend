package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/startdusk/finance-app-backend/internal/api/auth"
	v1 "github.com/startdusk/finance-app-backend/internal/api/v1"
	"github.com/startdusk/finance-app-backend/internal/database"
)

type API struct {
	Path            string
	Method          string
	Func            http.HandlerFunc
	permissionTypes []auth.PermissionType
}

func NewRouter(db database.Database) (http.Handler, error) {
	permissions := auth.NewPermissions(db)

	router := mux.NewRouter()
	router.HandleFunc("/version", v1.VersionHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	userAPI := &v1.UserAPI{
		DB: db,
	}

	accountAPI := &v1.AccountAPI{
		DB: db,
	}

	apis := []API{
		// ---------------USER-------------------
		NewAPI("/users", http.MethodPost, userAPI.Create, auth.Any),                             // Create user
		NewAPI("/users/{userID}", http.MethodGet, userAPI.Get, auth.Admin, auth.MemberIsTarget), // get user by id
		// NewAPI("/users/{userID}", http.MethodGet, userAPI.Get, auth.Admin, auth.MemberIsTarget), // list all user
		// NewAPI("/users/{userID}", http.MethodGet, userAPI.Get, auth.Admin, auth.MemberIsTarget), // delete user by id
		NewAPI("/login", http.MethodPost, userAPI.Login, auth.Any), // Login user

		// ---------------TOKENS------------------
		NewAPI("/refresh", http.MethodPost, userAPI.RefreshToken, auth.Any), // Refresh token

		// ---------------ROLES-------------------
		NewAPI("/users/{userID}/roles", http.MethodPost, userAPI.GrantRole, auth.Admin),    // Create role
		NewAPI("/users/{userID}/roles", http.MethodDelete, userAPI.RevokeRole, auth.Admin), // Revoke role
		NewAPI("/users/{userID}/roles", http.MethodGet, userAPI.GetRoleList, auth.Admin),   // Get all role

		// ---------------ACCOUNTS----------------
		NewAPI("/users/{userID}/accounts", http.MethodPost, accountAPI.Create, auth.Admin, auth.MemberIsTarget),               // create account for user (Open for admin for now)
		NewAPI("/users/{userID}/accounts", http.MethodGet, accountAPI.List, auth.Admin, auth.MemberIsTarget),                  // get accounts for user (Open for admin for now)
		NewAPI("/users/{userID}/accounts/{accountID}", http.MethodPatch, accountAPI.Update, auth.Admin, auth.MemberIsTarget),  // update account for user (Open for admin for now)
		NewAPI("/users/{userID}/accounts/{accountID}", http.MethodGet, accountAPI.Get, auth.Admin, auth.MemberIsTarget),       // get account by account id for user (Open for admin for now)
		NewAPI("/users/{userID}/accounts/{accountID}", http.MethodDelete, accountAPI.Delete, auth.Admin, auth.MemberIsTarget), // delete account by account id for user (Open for admin for now)
	}

	for _, api := range apis {
		apiRouter.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}

	router.Use(auth.AutherizationToken)

	return router, nil
}

func NewAPI(path string, method string, handlerFunc http.HandlerFunc, permissionTypes ...auth.PermissionType) API {
	return API{
		Path:            path,
		Method:          method,
		Func:            handlerFunc,
		permissionTypes: permissionTypes,
	}
}

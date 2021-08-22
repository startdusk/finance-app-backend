package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/startdusk/finance-app-backend/internal/api/auth"
	"github.com/startdusk/finance-app-backend/internal/api/utils"
	"github.com/startdusk/finance-app-backend/internal/database"
	"github.com/startdusk/finance-app-backend/internal/model"
)

// UserRoleAPI - providers REST for user role
type UserRoleAPI struct {
	DB database.Database
}

func SetUserRoleAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := &UserRoleAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{userID}/roles", api.GrantRole, auth.Admin),    // Create role
		NewAPI(http.MethodDelete, "/users/{userID}/roles", api.RevokeRole, auth.Admin), // Revoke role
		NewAPI(http.MethodGet, "/users/{userID}/roles", api.GetRoleList, auth.Admin),   // Get all role
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

func (api *UserRoleAPI) GrantRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go -> GrantRole()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	// Decode paramters
	var userRole model.UserRole
	if err := json.NewDecoder(r.Body).Decode(&userRole); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	// store role in database
	if err := api.DB.GrantRole(ctx, userID, userRole.Role); err != nil {
		logger.WithError(err).Warn("error granting role")
		utils.WriteError(w, http.StatusInternalServerError, "error granting role", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, &ActCreated{
		Created: true,
	})
}

func (api *UserRoleAPI) RevokeRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go -> RevokeRole()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	// Decode paramters
	var userRole model.UserRole
	if err := json.NewDecoder(r.Body).Decode(&userRole); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	// delete role from database
	if err := api.DB.RevokeRole(ctx, userID, userRole.Role); err != nil {
		logger.WithError(err).Warn("error revokting role")
		utils.WriteError(w, http.StatusInternalServerError, "error revokting role", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, &ActDeleted{
		Deleted: true,
	})
}

func (api *UserRoleAPI) GetRoleList(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go -> GetRoleList()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()

	// get user roles from database
	userRoles, err := api.DB.GetRolesByUser(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("error getting user roles")
		utils.WriteError(w, http.StatusInternalServerError, "error getting user roles", nil)
		return
	}
	if userRoles == nil {
		userRoles = make([]*model.UserRole, 0)
	}

	utils.WriteJSON(w, http.StatusCreated, &userRoles)
}

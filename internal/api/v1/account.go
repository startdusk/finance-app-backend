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

// AccountAPI - provides REST for Account
type AccountAPI struct {
	DB database.Database // will represent all database interface
}

func SetAccountAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := &AccountAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{userID}/accounts", api.Create, auth.Admin, auth.MemberIsTarget),               // create account for user (Open for admin for now)
		NewAPI(http.MethodGet, "/users/{userID}/accounts", api.List, auth.Admin, auth.MemberIsTarget),                  // get account for user (Open for admin for now)
		NewAPI(http.MethodPatch, "/users/{userID}/accounts/{accountID}", api.Update, auth.Admin, auth.MemberIsTarget),  // update account for user (Open for admin for now)
		NewAPI(http.MethodGet, "/users/{userID}/accounts/{accountID}", api.Get, auth.Admin, auth.MemberIsTarget),       // get account by account id for user (Open for admin for now)
		NewAPI(http.MethodDelete, "/users/{userID}/accounts/{accountID}", api.Delete, auth.Admin, auth.MemberIsTarget), // delete account by account id for user (Open for admin for now)
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

// POST - /users/{userID}/accounts
// Permission - MemberIsTarget
func (api *AccountAPI) Create(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "account.go -> Create()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	// Decode paramters
	var account model.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	account.UserID = &userID

	if err := account.Verify(); err != nil {
		logger.WithError(err).Warn("not all fields found") // I will hide this error in future, it isn't secure to show what fields are missing...
		utils.WriteError(w, http.StatusBadRequest, "not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.CreateAccount(ctx, &account); err != nil {
		logger.WithError(err).Warn("error creating account")
		utils.WriteError(w, http.StatusInternalServerError, "error creating account", nil)
		return
	}

	logger.WithField("accountID", account.ID).Info("account created")

	utils.WriteJSON(w, http.StatusCreated, &account)
}

// PATCH - /users/{userID}/accounts/{accountID}
// Permission - MemberIsTarget
func (api *AccountAPI) Update(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "account.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	accountID := model.AccountID(vars["accountID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
		"accountID": accountID,
	})

	// Decode paramters
	var accountRequest model.Account
	if err := json.NewDecoder(r.Body).Decode(&accountRequest); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithField("accountID", accountID)

	ctx := r.Context()

	account, err := api.DB.GetAccountByID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Warn("error getting account")
		utils.WriteError(w, http.StatusConflict, "error getting account", nil)
		return
	}

	if accountRequest.Name != nil || len(*accountRequest.Name) != 0 {
		account.Name = accountRequest.Name
	}

	if accountRequest.Type != nil || len(*accountRequest.Type) != 0 {
		account.Type = accountRequest.Type
	}

	if accountRequest.StartBalance != nil {
		account.StartBalance = accountRequest.StartBalance
	}

	if accountRequest.Currency != nil || len(*accountRequest.Currency) != 0 {
		account.Currency = accountRequest.Currency
	}

	if err := api.DB.UpdateAccount(ctx, account); err != nil {
		logger.WithError(err).Warn("error updating account")
		utils.WriteError(w, http.StatusInternalServerError, "error updating account", nil)
		return
	}

	logger.Info("account updated")

	utils.WriteJSON(w, http.StatusOK, &ActUpdated{
		Updated: true,
	})
}

// GET - /users/{userID}/accounts
// Permission - MemberIsTarget
func (api *AccountAPI) List(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "account.go -> List()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()

	accounts, err := api.DB.ListAccountsByUserID(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("error getting accounts")
		utils.WriteError(w, http.StatusConflict, "error getting accounts", nil)
		return
	}

	logger.Info("accounts returned")

	utils.WriteJSON(w, http.StatusOK, &accounts)
}

// GET - /users/{userID}/accounts/{accountID}
// Permission - MemberIsTarget
func (api *AccountAPI) Get(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "account.go -> Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	accountID := model.AccountID(vars["accountID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
		"accountID": accountID,
	})

	ctx := r.Context()

	account, err := api.DB.GetAccountByID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Warn("error getting account")
		utils.WriteError(w, http.StatusConflict, "error getting account", nil)
		return
	}

	logger.Info("account returned")

	utils.WriteJSON(w, http.StatusOK, &account)
}

// DELETE - /users/{userID}/accounts/{accountID}
// Permission - MemberIsTarget
func (api *AccountAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "account.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	accountID := model.AccountID(vars["accountID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
		"accountID": accountID,
	})

	ctx := r.Context()

	ok, err := api.DB.DeleteAccount(ctx, accountID)
	if !ok && err != nil {
		logger.WithError(err).Warn("error deleting account")
		utils.WriteError(w, http.StatusConflict, "error deleting account", nil)
		return
	}

	logger.Info("account deleted")

	utils.WriteJSON(w, http.StatusOK, &ActDeleted{
		Deleted: true,
	})
}

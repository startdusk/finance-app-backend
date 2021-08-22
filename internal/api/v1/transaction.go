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

// TransactionAPI - provides REST for Transaction
type TransactionAPI struct {
	DB database.Database // will represent all database interface
}

func SetTransactionAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := &TransactionAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{userID}/transactions", api.Create, auth.Admin, auth.MemberIsTarget),                   // create transaction for user (Open for admin for now)
		NewAPI(http.MethodGet, "/users/{userID}/transactions", api.ListByUser, auth.Admin, auth.MemberIsTarget),                // get transaction for user (Open for admin for now)
		NewAPI(http.MethodGet, "/accounts/{accountID}/transactions", api.ListByAccount, auth.Admin, auth.MemberIsTarget),       // get transaction for account (Open for admin for now)
		NewAPI(http.MethodGet, "/categories/{categoryID}/transactions", api.ListByCategory, auth.Admin, auth.MemberIsTarget),   // get transaction for category (Open for admin for now)
		NewAPI(http.MethodPatch, "/users/{userID}/transactions/{transactionID}", api.Update, auth.Admin, auth.MemberIsTarget),  // update transaction for user (Open for admin for now)
		NewAPI(http.MethodGet, "/users/{userID}/transactions/{transactionID}", api.Get, auth.Admin, auth.MemberIsTarget),       // get transaction by transaction id for user (Open for admin for now)
		NewAPI(http.MethodDelete, "/users/{userID}/transactions/{transactionID}", api.Delete, auth.Admin, auth.MemberIsTarget), // delete transaction by transaction id for user (Open for admin for now)
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

// POST - /users/{userID}/transactions
// Permission - MemberIsTarget
func (api *TransactionAPI) Create(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "transaction.go -> Create()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	// Decode paramters
	var transaction model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	transaction.UserID = &userID

	if err := transaction.Verify(); err != nil {
		logger.WithError(err).Warn("not all fields found") // I will hide this error in future, it isn't secure to show what fields are missing...
		utils.WriteError(w, http.StatusBadRequest, "not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.CreateTransaction(ctx, &transaction); err != nil {
		logger.WithError(err).Warn("error creating transaction")
		utils.WriteError(w, http.StatusInternalServerError, "error creating transaction", nil)
		return
	}

	logger.WithField("TransactionID", transaction.ID).Info("transaction created")

	utils.WriteJSON(w, http.StatusCreated, &transaction)
}

// PATCH - /users/{userID}/transactions/{transactionID}
// Permission - MemberIsTarget
func (api *TransactionAPI) Update(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "transaction.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	transactionID := model.TransactionID(vars["transactionID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":        userID,
		"principal":     principal,
		"transactionID": transactionID,
	})

	// Decode paramters
	var transactionRequest model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transactionRequest); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithField("transactionID", transactionID)

	ctx := r.Context()

	transaction, err := api.DB.GetTransactionByID(ctx, transactionID)
	if err != nil {
		logger.WithError(err).Warn("error getting transaction")
		utils.WriteError(w, http.StatusConflict, "error getting transaction", nil)
		return
	}

	if transactionRequest.AccountID != nil || *transactionRequest.AccountID != model.NilAccountID {
		transaction.AccountID = transactionRequest.AccountID
	}

	if transactionRequest.CategoryID != nil || *transactionRequest.CategoryID != model.NilCategoryID {
		transaction.CategoryID = transactionRequest.CategoryID
	}

	if transactionRequest.Date != nil {
		transaction.Date = transactionRequest.Date
	}

	if transactionRequest.Type != nil || *transactionRequest.Type != "" {
		transaction.Type = transactionRequest.Type
	}

	if transactionRequest.Amount != nil {
		transaction.Amount = transactionRequest.Amount
	}

	if transactionRequest.Notes != nil {
		transaction.Notes = transactionRequest.Notes
	}

	if err := api.DB.UpdateTransaction(ctx, transaction); err != nil {
		logger.WithError(err).Warn("error updating transaction")
		utils.WriteError(w, http.StatusInternalServerError, "error updating transaction", nil)
		return
	}

	logger.Info("transaction updated")

	utils.WriteJSON(w, http.StatusOK, &ActUpdated{
		Updated: true,
	})
}

// GET - /users/{userID}/transactions?from={from}&to={to}
// Permission - MemberIsTarget
func (api *TransactionAPI) ListByUser(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "transaction.go -> ListByUser()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	query := r.URL.Query()
	from, err := utils.TimeParam(query, "from")
	if err != nil {
		logger.WithError(err).Warn("invalid from parameters")
		utils.WriteError(w, http.StatusConflict, "invalid from parameters", nil)
		return
	}

	to, err := utils.TimeParam(query, "to")
	if err != nil {
		logger.WithError(err).Warn("invalid to parameters")
		utils.WriteError(w, http.StatusConflict, "invalid to parameters", nil)
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
		"from":      from,
		"to":        to,
	})

	ctx := r.Context()

	transactions, err := api.DB.ListTransactionByUserID(ctx, userID, from, to)
	if err != nil {
		logger.WithError(err).Warn("error getting transactions")
		utils.WriteError(w, http.StatusConflict, "error getting transactions", nil)
		return
	}

	logger.Info("transactions returned")

	if transactions == nil {
		transactions = make([]*model.Transaction, 0)
	}

	utils.WriteJSON(w, http.StatusOK, &transactions)
}

// GET - /categories/{categoryID}/transactions?from={from}&to={to}
// Permission - MemberIsTarget
func (api *TransactionAPI) ListByCategory(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "transaction.go -> ListByCategory()")

	vars := mux.Vars(r)
	categoryID := model.CategoryID(vars["categoryID"])
	principal := auth.GetPrincipal(r)

	query := r.URL.Query()
	from, err := utils.TimeParam(query, "from")
	if err != nil {
		logger.WithError(err).Warn("invalid from parameters")
		utils.WriteError(w, http.StatusConflict, "invalid from parameters", nil)
		return
	}

	to, err := utils.TimeParam(query, "to")
	if err != nil {
		logger.WithError(err).Warn("invalid to parameters")
		utils.WriteError(w, http.StatusConflict, "invalid to parameters", nil)
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"categoryID": categoryID,
		"principal":  principal,
		"from":       from,
		"to":         to,
	})

	ctx := r.Context()

	transactions, err := api.DB.ListTransactionByCategoryID(ctx, categoryID, from, to)
	if err != nil {
		logger.WithError(err).Warn("error getting transactions")
		utils.WriteError(w, http.StatusConflict, "error getting transactions", nil)
		return
	}

	logger.Info("transactions returned")

	if transactions == nil {
		transactions = make([]*model.Transaction, 0)
	}

	utils.WriteJSON(w, http.StatusOK, &transactions)
}

// GET - /accounts/{accountID}/transactions?from={from}&to={to}
// Permission - MemberIsTarget
func (api *TransactionAPI) ListByAccount(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "transaction.go -> ListByAccount()")

	vars := mux.Vars(r)
	accountID := model.AccountID(vars["accountID"])
	principal := auth.GetPrincipal(r)

	query := r.URL.Query()
	from, err := utils.TimeParam(query, "from")
	if err != nil {
		logger.WithError(err).Warn("invalid from parameters")
		utils.WriteError(w, http.StatusConflict, "invalid from parameters", nil)
		return
	}

	to, err := utils.TimeParam(query, "to")
	if err != nil {
		logger.WithError(err).Warn("invalid to parameters")
		utils.WriteError(w, http.StatusConflict, "invalid to parameters", nil)
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"accountID": accountID,
		"principal": principal,
		"from":      from,
		"to":        to,
	})

	ctx := r.Context()

	transactions, err := api.DB.ListTransactionByAccountID(ctx, accountID, from, to)
	if err != nil {
		logger.WithError(err).Warn("error getting transactions")
		utils.WriteError(w, http.StatusConflict, "error getting transactions", nil)
		return
	}

	logger.Info("transactions returned")

	if transactions == nil {
		transactions = make([]*model.Transaction, 0)
	}

	utils.WriteJSON(w, http.StatusOK, &transactions)
}

// GET - /users/{userID}/transactions/{transactionID}
// Permission - MemberIsTarget
func (api *TransactionAPI) Get(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "transaction.go -> Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	transactionID := model.TransactionID(vars["transactionID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":        userID,
		"principal":     principal,
		"transactionID": transactionID,
	})

	ctx := r.Context()

	transaction, err := api.DB.GetTransactionByID(ctx, transactionID)
	if err != nil {
		logger.WithError(err).Warn("error getting transaction")
		utils.WriteError(w, http.StatusConflict, "error getting transaction", nil)
		return
	}

	logger.Info("transaction returned")

	utils.WriteJSON(w, http.StatusOK, &transaction)
}

// DELETE - /users/{userID}/transactions/{transactionID}
// Permission - MemberIsTarget
func (api *TransactionAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "transaction.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	transactionID := model.TransactionID(vars["transactionID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":        userID,
		"principal":     principal,
		"transactionID": transactionID,
	})

	ctx := r.Context()

	ok, err := api.DB.DeleteTransaction(ctx, transactionID)
	if !ok && err != nil {
		logger.WithError(err).Warn("error deleting transaction")
		utils.WriteError(w, http.StatusConflict, "error deleting transaction", nil)
		return
	}

	logger.Info("transaction deleted")

	utils.WriteJSON(w, http.StatusOK, &ActDeleted{
		Deleted: true,
	})
}

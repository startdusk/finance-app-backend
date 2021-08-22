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

// MerchantAPI - provides REST for Merchant
type MerchantAPI struct {
	DB database.Database // will represent all database interface
}

func SetMerchantAPI(db database.Database, router *mux.Router, permissions auth.Permissions) {
	api := &MerchantAPI{
		DB: db,
	}

	apis := []API{
		NewAPI(http.MethodPost, "/users/{userID}/merchants", api.Create, auth.Admin, auth.MemberIsTarget),                // create merchant for user (Open for admin for now)
		NewAPI(http.MethodGet, "/users/{userID}/merchants", api.List, auth.Admin, auth.MemberIsTarget),                   // get merchant for user (Open for admin for now)
		NewAPI(http.MethodPatch, "/users/{userID}/merchants/{merchantID}", api.Update, auth.Admin, auth.MemberIsTarget),  // update merchant for user (Open for admin for now)
		NewAPI(http.MethodGet, "/users/{userID}/merchants/{merchantID}", api.Get, auth.Admin, auth.MemberIsTarget),       // get merchant by merchant id for user (Open for admin for now)
		NewAPI(http.MethodDelete, "/users/{userID}/merchants/{merchantID}", api.Delete, auth.Admin, auth.MemberIsTarget), // delete merchant by merchant id for user (Open for admin for now)
	}

	for _, api := range apis {
		router.HandleFunc(api.Path, permissions.Wrap(api.Func, api.permissionTypes...)).Methods(api.Method)
	}
}

// POST - /users/{userID}/merchants
// Permission - MemberIsTarget
func (api *MerchantAPI) Create(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "merchant.go -> Create()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	// Decode paramters
	var merchant model.Merchant
	if err := json.NewDecoder(r.Body).Decode(&merchant); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	merchant.UserID = &userID

	if err := merchant.Verify(); err != nil {
		logger.WithError(err).Warn("not all fields found") // I will hide this error in future, it isn't secure to show what fields are missing...
		utils.WriteError(w, http.StatusBadRequest, "not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.CreateMerchant(ctx, &merchant); err != nil {
		logger.WithError(err).Warn("error creating merchant")
		utils.WriteError(w, http.StatusInternalServerError, "error creating merchant", nil)
		return
	}

	logger.WithField("MerchantID", merchant.ID).Info("merchant created")

	utils.WriteJSON(w, http.StatusCreated, &merchant)
}

// PATCH - /users/{userID}/merchants/{merchantID}
// Permission - MemberIsTarget
func (api *MerchantAPI) Update(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "merchant.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	merchantID := model.MerchantID(vars["merchantID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"merchantID": merchantID,
	})

	// Decode paramters
	var merchantRequest model.Merchant
	if err := json.NewDecoder(r.Body).Decode(&merchantRequest); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithField("merchantID", merchantID)

	ctx := r.Context()

	merchant, err := api.DB.GetMerchantByID(ctx, merchantID)
	if err != nil {
		logger.WithError(err).Warn("error getting merchant")
		utils.WriteError(w, http.StatusConflict, "error getting merchant", nil)
		return
	}

	if merchantRequest.Name != nil || len(*merchantRequest.Name) != 0 {
		merchant.Name = merchantRequest.Name
	}

	if err := api.DB.UpdateMerchant(ctx, merchant); err != nil {
		logger.WithError(err).Warn("error updating merchant")
		utils.WriteError(w, http.StatusInternalServerError, "error updating merchant", nil)
		return
	}

	logger.Info("merchant updated")

	utils.WriteJSON(w, http.StatusOK, &ActUpdated{
		Updated: true,
	})
}

// GET - /users/{userID}/merchants
// Permission - MemberIsTarget
func (api *MerchantAPI) List(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "merchant.go -> List()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()

	merchants, err := api.DB.ListMerchantsByUserID(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("error getting merchants")
		utils.WriteError(w, http.StatusConflict, "error getting merchants", nil)
		return
	}

	logger.Info("merchants returned")

	utils.WriteJSON(w, http.StatusOK, &merchants)
}

// GET - /users/{userID}/merchants/{MerchantID}
// Permission - MemberIsTarget
func (api *MerchantAPI) Get(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "merchant.go -> Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	merchantID := model.MerchantID(vars["merchantID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"merchantID": merchantID,
	})

	ctx := r.Context()

	merchant, err := api.DB.GetMerchantByID(ctx, merchantID)
	if err != nil {
		logger.WithError(err).Warn("error getting merchant")
		utils.WriteError(w, http.StatusConflict, "error getting merchant", nil)
		return
	}

	logger.Info("merchant returned")

	utils.WriteJSON(w, http.StatusOK, &merchant)
}

// DELETE - /users/{userID}/merchants/{merchantID}
// Permission - MemberIsTarget
func (api *MerchantAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "merchant.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	merchantID := model.MerchantID(vars["merchantID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"merchantID": merchantID,
	})

	ctx := r.Context()

	ok, err := api.DB.DeleteMerchant(ctx, merchantID)
	if !ok && err != nil {
		logger.WithError(err).Warn("error deleting merchant")
		utils.WriteError(w, http.StatusConflict, "error deleting merchant", nil)
		return
	}

	logger.Info("merchant deleted")

	utils.WriteJSON(w, http.StatusOK, &ActDeleted{
		Deleted: true,
	})
}

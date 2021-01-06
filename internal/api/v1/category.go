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

// CategoryAPI - provides REST for Category
type CategoryAPI struct {
	DB database.Database // will represent all database interface
}

// POST - /users/{userID}/categories
// Permission - MemberIsTarget
func (api *CategoryAPI) Create(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "category.go -> Create()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	// Decode paramters
	var category model.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	category.UserID = &userID

	if err := category.Verify(); err != nil {
		logger.WithError(err).Warn("not all fields found") // I will hide this error in future, it isn't secure to show what fields are missing...
		utils.WriteError(w, http.StatusBadRequest, "not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.CreateCategory(ctx, &category); err != nil {
		logger.WithError(err).Warn("error creating category")
		utils.WriteError(w, http.StatusInternalServerError, "error creating category", nil)
		return
	}

	logger.WithField("categoryID", category.ID).Info("category created")

	utils.WriteJSON(w, http.StatusCreated, &category)
}

// PATCH - /users/{userID}/categories/{categoryID}
// Permission - MemberIsTarget
func (api *CategoryAPI) Update(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "category.go -> Update()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	categoryID := model.CategoryID(vars["categoryID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"categoryID": categoryID,
	})

	// Decode paramters
	var categoryRequest model.Category
	if err := json.NewDecoder(r.Body).Decode(&categoryRequest); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithField("categoryID", categoryID)

	ctx := r.Context()

	category, err := api.DB.GetCategoryByID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Warn("error getting category")
		utils.WriteError(w, http.StatusConflict, "error getting category", nil)
		return
	}

	if categoryRequest.ParentID != model.NilCategoryID {
		category.ParentID = categoryRequest.ParentID
	}

	if categoryRequest.Name != nil || len(*categoryRequest.Name) != 0 {
		category.Name = categoryRequest.Name
	}

	if err := api.DB.UpdateCategory(ctx, category); err != nil {
		logger.WithError(err).Warn("error updating category")
		utils.WriteError(w, http.StatusInternalServerError, "error updating category", nil)
		return
	}

	logger.Info("category updated")

	utils.WriteJSON(w, http.StatusOK, &ActUpdated{
		Updated: true,
	})
}

// GET - /users/{userID}/categories
// Permission - MemberIsTarget
func (api *CategoryAPI) List(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "category.go -> List()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()

	categories, err := api.DB.ListCategoriesByUserID(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("error getting categories")
		utils.WriteError(w, http.StatusConflict, "error getting categories", nil)
		return
	}

	logger.Info("categories returned")

	utils.WriteJSON(w, http.StatusOK, &categories)
}

// GET - /users/{userID}/categories/{categoryID}
// Permission - MemberIsTarget
func (api *CategoryAPI) Get(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "category.go -> Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	categoryID := model.CategoryID(vars["categoryID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"categoryID": categoryID,
	})

	ctx := r.Context()

	category, err := api.DB.GetCategoryByID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Warn("error getting category")
		utils.WriteError(w, http.StatusConflict, "error getting category", nil)
		return
	}

	logger.Info("category returned")

	utils.WriteJSON(w, http.StatusOK, &category)
}

// DELETE - /users/{userID}/categories/{categoryID}
// Permission - MemberIsTarget
func (api *CategoryAPI) Delete(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "category.go -> Delete()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	categoryID := model.CategoryID(vars["categoryID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":     userID,
		"principal":  principal,
		"categoryID": categoryID,
	})

	ctx := r.Context()

	ok, err := api.DB.DeleteCategory(ctx, categoryID)
	if !ok && err != nil {
		logger.WithError(err).Warn("error deleting category")
		utils.WriteError(w, http.StatusConflict, "error deleting category", nil)
		return
	}

	logger.Info("category deleted")

	utils.WriteJSON(w, http.StatusOK, &ActDeleted{
		Deleted: true,
	})
}

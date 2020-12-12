package v1

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/startdusk/finance-app-backend/internal/api/utils"
	"github.com/startdusk/finance-app-backend/internal/database"
	"github.com/startdusk/finance-app-backend/internal/model"
)

type UserAPI struct {
	DB database.Database // will represent all database interface
}

type UserParameters struct {
	model.User
	Password string `json:"password"`
}

func (api *UserAPI) Create(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "user.go -> Create()")
	// Load parameters
	var userParameters UserParameters
	if err := json.NewDecoder(r.Body).Decode(&userParameters); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"email": *userParameters.Email,
	})

	if err := userParameters.Verify(); err != nil {
		logger.WithError(err).Warn("not all fields found")
		utils.WriteError(w, http.StatusBadRequest, "not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	hashed, err := model.HashPassword(userParameters.Password)
	if err != nil {
		logger.WithError(err).Warn("could not hash password")
		utils.WriteError(w, http.StatusInternalServerError, "could not hash password", nil)
		return
	}

	newUser := &model.User{
		Email:        userParameters.Email,
		PasswordHash: &hashed,
	}

	ctx := r.Context()
	if err := api.DB.CreateUser(ctx, newUser); err == database.ErrUserExist {
		logger.WithError(err).Warn("user already exists")
		utils.WriteError(w, http.StatusConflict, "user already exists", nil)
		return
	} else if err != nil {
		logger.WithError(err).Warn("error creating user")
		utils.WriteError(w, http.StatusConflict, "error creating user", nil)
		return
	}

	createdUser, err := api.DB.GetUserByID(ctx, &newUser.ID)
	if err != nil {
		logger.WithError(err).Warn("error creating user")
		utils.WriteError(w, http.StatusConflict, "error creating user", nil)
		return
	}

	logger.Info("user created")

	utils.WriteJSON(w, http.StatusCreated, createdUser)
}

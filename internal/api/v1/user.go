package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/startdusk/finance-app-backend/internal/api/auth"
	"github.com/startdusk/finance-app-backend/internal/api/utils"
	"github.com/startdusk/finance-app-backend/internal/database"
	"github.com/startdusk/finance-app-backend/internal/model"
)

// UserAPI - providers REST for users
type UserAPI struct {
	DB database.Database // will represent all database interface
}

type UserParameters struct {
	model.User
	model.SessionData

	Password string `json:"password"` // Password must be 8 characters or longer!
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

	if err := userParameters.User.Verify(); err != nil {
		logger.WithError(err).Warn("not all fields found") // I will hide this error in future, it isn't secure to show what fields are missing...
		utils.WriteError(w, http.StatusBadRequest, "not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := userParameters.SessionData.Verify(); err != nil {
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

	logger.WithField("userID", createdUser.ID).Info("user created")

	api.writeTokenResponse(ctx, w, http.StatusCreated, createdUser, &userParameters.SessionData, true)
}

func (api *UserAPI) Login(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "user.go -> Login()")

	var credantials model.Credantials
	if err := json.NewDecoder(r.Body).Decode(&credantials); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := credantials.SessionData.Verify(); err != nil {
		logger.WithError(err).Warn("not all fields found")
		utils.WriteError(w, http.StatusBadRequest, "not all fields found", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"email": credantials.Email,
	})

	ctx := r.Context()
	user, err := api.DB.GetUserByEmail(ctx, credantials.Email)
	if err != nil {
		logger.WithError(err).Warn("error login in")
		utils.WriteError(w, http.StatusConflict, "invalid email or password", nil)
		return
	}

	// Checking if password is correct
	if err := user.CheckPassword(credantials.Password); err != nil {
		logger.WithError(err).Warn("error login in")
		utils.WriteError(w, http.StatusConflict, "invalid email or password", nil)
		return
	}

	logger.WithField("userID", user.ID).Info("user login in")

	api.writeTokenResponse(ctx, w, http.StatusOK, user, &credantials.SessionData, true)
}

func (api *UserAPI) Get(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "user.go -> Get()")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])

	ctx := r.Context()
	user, err := api.DB.GetUserByID(ctx, &userID)
	if err != nil {
		logger.WithError(err).Warn("error getting user")
		utils.WriteError(w, http.StatusConflict, "error getting user", nil)
		return
	}

	logger.WithField("userID", user.ID).Info("get user complete")

	utils.WriteJSON(w, http.StatusOK, user)
}

// RefreshTokenRequest - Data user sned to get new access refresh tokens
type RefreshTokenRequest struct {
	RefreshToken string         `json:"refreshToken"`
	DeviceID     model.DeviceID `json:"deviceID"`
}

func (api *UserAPI) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// show function name to track error faster
	logger := logrus.WithField("func", "user.go -> RefreshToken()")

	// TODO: move this part to separate function
	var request RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"DeviceID": request.DeviceID,
	})

	principal, err := auth.VerifyToken(request.RefreshToken)
	if err != nil {
		logger.WithError(err).Warn("error verifing refresh token")
		utils.WriteError(w, http.StatusUnauthorized, "error verifing refresh token", nil)
		return
	}

	// if token is valid we need to check if combination UserID - DeviceID Refresh Token exists in database
	data := model.Session{
		UserID:       principal.UserID,
		DeviceID:     request.DeviceID,
		RefreshToken: request.RefreshToken,
	}

	ctx := r.Context()
	session, err := api.DB.GetSession(ctx, data)
	if err != nil || session == nil {
		logger.WithError(err).Warn("error session not exists")
		utils.WriteError(w, http.StatusUnauthorized, "error session not exists", nil)
		return
	}

	// if session exists and valid we generate new access and refresh tokens.
	logger.WithField("userID", principal.UserID).Debug("refresh token")

	// check if user exists
	user, err := api.DB.GetUserByID(ctx, &principal.UserID)
	if err != nil {
		logger.WithError(err).Warn("error getting user")
		utils.WriteError(w, http.StatusConflict, "error getting user", nil)
		return
	}
	api.writeTokenResponse(ctx, w, http.StatusOK, user, &model.SessionData{DeviceID: request.DeviceID}, true)
}

type TokenResponse struct {
	Tokens *auth.Tokens `json:"tokens,omitempty"` // this will insert all tokens struct fields
	User   *model.User  `json:"user,omitempty"`
}

// writeTokenResponse - Generate Access and Refresh token are return them to user. Refresh token is stored in database as session
func (api *UserAPI) writeTokenResponse(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	user *model.User,
	sessionData *model.SessionData,
	cookie bool) {
	// Issue token:
	// TODO: add user role to Principal
	tokens, err := auth.IssueToken(model.Principal{UserID: user.ID})
	if err != nil && tokens == nil {
		logrus.WithError(err).Warn("error issuing token")
		utils.WriteError(w, http.StatusConflict, "error issuing token", nil)
		return
	}

	session := model.Session{
		UserID:       user.ID,
		DeviceID:     sessionData.DeviceID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.RefreshTokenExpiresAt, // 存储的refreshToken的过期时间，返回前端的是accessToken的时间
	}

	if err := api.DB.SaveRefreshToken(ctx, session); err != nil {
		logrus.WithError(err).Warn("error issuing token")
		utils.WriteError(w, http.StatusConflict, "error issuing token", nil)
		return
	}

	// Write token response
	tokenResponse := TokenResponse{
		Tokens: tokens,
		User:   user,
	}

	if cookie {
		// later
	}

	utils.WriteJSON(w, status, tokenResponse)
}

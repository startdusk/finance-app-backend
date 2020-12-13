package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/startdusk/finance-app-backend/internal/api/auth"
	v1 "github.com/startdusk/finance-app-backend/internal/api/v1"
	"github.com/startdusk/finance-app-backend/internal/database"
)

func NewRouter(db database.Database, tokens auth.Tokens) (http.Handler, error) {
	router := mux.NewRouter()
	router.HandleFunc("/version", v1.VersionHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	userAPI := &v1.UserAPI{
		DB:     db,
		Tokens: tokens,
	}
	apiRouter.HandleFunc("/users", userAPI.Create).Methods(http.MethodPost)
	//apiRouter.HandleFunc("/users", userAPI.Create).Methods(http.MethodGet) // list all user
	//apiRouter.HandleFunc("/users/{userID}", userAPI.Create).Methods(http.MethodGet) // get user by id
	//apiRouter.HandleFunc("/users/{userID}", userAPI.Create).Methods(http.MethodDelete) // delete user by id

	apiRouter.HandleFunc("/login", userAPI.Login).Methods(http.MethodPost)

	return router, nil
}

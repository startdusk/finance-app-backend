package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/startdusk/finance-app-backend/internal/api/auth"
	v1 "github.com/startdusk/finance-app-backend/internal/api/v1"
	"github.com/startdusk/finance-app-backend/internal/database"
)

func NewRouter(db database.Database) (http.Handler, error) {
	permissions := auth.NewPermissions(db)

	router := mux.NewRouter()
	router.HandleFunc("/version", v1.VersionHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	v1.SetUserAPI(db, apiRouter, permissions)
	v1.SetUserRoleAPI(db, apiRouter, permissions)
	v1.SetAccountAPI(db, apiRouter, permissions)
	v1.SetCategoryAPI(db, apiRouter, permissions)
	v1.SetMerchantAPI(db, apiRouter, permissions)
	v1.SetTransactionAPI(db, apiRouter, permissions)
	router.Use(auth.AutherizationToken)

	return router, nil
}

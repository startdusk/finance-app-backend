package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/startdusk/finance-app-backend/internal/api/v1"
	"github.com/startdusk/finance-app-backend/internal/database"
)

func NewRouter(db database.Database) (http.Handler, error) {
	router := mux.NewRouter()
	router.HandleFunc("/version", v1.VersionHandler)
	return router, nil
}
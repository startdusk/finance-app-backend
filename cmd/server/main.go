package main

import (
	"net/http"

	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"

	"github.com/startdusk/finance-app-backend/internal/api"
	"github.com/startdusk/finance-app-backend/internal/api/auth"
	"github.com/startdusk/finance-app-backend/internal/config"
	"github.com/startdusk/finance-app-backend/internal/database"
)

// Create Server object and listener
func main() {
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithField("version", config.Version).Debug("Starting server.")

	tokens := auth.NewTokens()

	// Creating new database
	db, err := database.New()
	if err != nil {
		logrus.WithError(err).Fatal("Error verifying database")
	}

	router, err := api.NewRouter(db, tokens)
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}

	logrus.Debug("Database is ready to use")

	const addr = "0.0.0.0:8088"
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	// Starting server
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed.")
	}
}

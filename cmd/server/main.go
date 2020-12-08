package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/startdusk/finance-app-backend/internal/api"
	"github.com/startdusk/finance-app-backend/internal/config"
)

// Create Server object and listener
func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithField("version", config.Version).Debug("Starting server.")

	router, err := api.NewRouter()
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}

	const addr = "0.0.0.0:8088"
	server := http.Server{
		Handler: router,
		Addr: addr,
	}

	// Starting server
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed")
	}
}

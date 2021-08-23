package main

import (
	"net"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"

	"github.com/startdusk/finance-app-backend/internal/api"
	"github.com/startdusk/finance-app-backend/internal/config"
	"github.com/startdusk/finance-app-backend/internal/database"
)

var (
	host = flag.String("host", "0.0.0.0", "host for listen")
	port = flag.String("port", "8088", "port for listen")
)

// Create Server object and listener
func main() {
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithField("version", config.Version).Debug("Starting server.")

	// Creating new database
	db, err := database.New()
	if err != nil {
		logrus.WithError(err).Fatal("Error verifying database")
	}

	router, err := api.NewRouter(db)
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}

	logrus.Debug("Database is ready to use")

	var addr = net.JoinHostPort(*host, *port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	// Starting server
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed.")
	}
}

package database

import (
	"time"
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/namsral/flag"
	"github.com/jmoiron/sqlx"
)

var (
	databaseURL = flag.String("database-url", "postgres://postgres:password@db:5432/postgres?sslmode=disable", "Database URL")
	databaseTimeout = flag.Int64("database-timeout-ms", 2000, "")
)

// Connect creates a new database connection
func Connect() (*sqlx.DB, error) {
	// Connect to database
	dbURL := *databaseURL

	logrus.WithField("url", dbURL).Debug("Connecting to database")
	conn, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to database")
	}

	conn.SetMaxOpenConns(32)

	// Check if database running
	if err := waitForDB(conn.DB); err != nil {
		return nil, err
	}

	// Migrate database schema
	if err := migrateDB(conn.DB); err != nil {
		return nil, errors.Wrap(err, "could not migrate database")
	}

	return conn, nil
}

// New creates a new database
func New() (Database, error) {
	conn, err := Connect()
	if err != nil {
		return nil, err
	}

	d := &database{
		conn: conn,
	}
	return d, nil
}


func waitForDB(conn *sql.DB) error {
	ready := make(chan struct{})
	go func() {
		for {
			if err := conn.Ping(); err == nil {
				close(ready)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		select {
		case <-ready:
			return nil
		case <-time.After(time.Duration(*databaseTimeout) * time.Millisecond):
			return errors.New("database not ready")
		}
	}
}

package database

import (
	"io"

	"github.com/jmoiron/sqlx"
)

// Database - interface for database
type Database interface {
	io.Closer
}

type database struct {
	conn *sqlx.DB
}

func (d *database) Close() error {
	return d.conn.Close()
}
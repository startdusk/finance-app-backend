package database

import (
	"io"

	"github.com/jmoiron/sqlx"
)

// UniqueViolation Postgres error string for unique index violation(unique index violation: 唯一索引冲突)
const UniqueViolation = "unique_violation"

// Database - interface for database
type Database interface {
	UsersDB
	SessionDB
	UserRoleDB
	AccountDB
	CategoryDB

	io.Closer
}

type database struct {
	conn *sqlx.DB
}

func (d *database) Close() error {
	return d.conn.Close()
}

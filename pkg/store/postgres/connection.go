package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// NewConnectionInput input for NewConnection
type NewConnectionInput struct {
	User     string
	Password string
	Dbname   string
	Host     string
	Port     string
}

// NewConnection Helper for creating postgres connection
func NewConnection(i *NewConnectionInput) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", i.User, i.Password, i.Host, i.Port, i.Dbname)
	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "Database connetion failed")
	}

	err = connection.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}
	connection.SetMaxIdleConns(0) // Let Aurora sleep if no connections present.

	return connection, nil
}

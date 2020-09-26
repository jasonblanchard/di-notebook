package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://pkg/db/migrations",
		"postgres://di:di@localhost:5432/di_notebook?sslmode=disable")

	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		panic(err)
	}
}

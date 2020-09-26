package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Up | Down")
		return
	}

	cmd := os.Args[1]
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")

	m, err := migrate.New(
		"file://pkg/db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database))

	if err != nil {
		panic(err)
	}

	switch cmd {
	case "up":
		err = m.Up()
		if err != nil {
			panic(err)
		}
		log.Print("Finished UP")
	case "down":
		err = m.Down()
		if err != nil {
			panic(err)
		}
		log.Print("Finished DOWN")
	}

}

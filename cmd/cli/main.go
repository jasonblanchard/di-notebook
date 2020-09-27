package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jasonblanchard/di-notebook/cmd/cli/cmd"
)

func main() {
	cmd.Execute()
}

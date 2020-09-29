package app

import (
	"testing"

	"github.com/jasonblanchard/di-notebook/store/postgres"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties/assert"
)

func makeApp() (*App, error) {
	db, err := postgres.NewConnection(&postgres.NewConnectionInput{
		User:     "di",
		Password: "di",
		Dbname:   "di_notebook",
		Host:     "localhost",
	})

	if err != nil {
		return nil, err
	}

	reader := &postgres.Reader{
		Db: db,
	}

	writer := &postgres.Writer{
		Db: db,
	}

	app := &App{
		StoreReader: reader,
		StoreWriter: writer,
	}

	return app, nil
}

func TestCreateReadFlow(t *testing.T) {
	app, err := makeApp()

	if err != nil {
		panic(err)
	}

	tester := &Principle{
		Type: PrincipleTypeTest,
	}

	err = app.ResetEntries(tester)

	if err != nil {
		panic(err)
	}

	author := &Principle{
		Type: PrincipleTypeUser,
		ID:   "123",
	}

	id, err := app.StartNewEntry(author, "hello", "123")
	if err != nil {
		panic(err)
	}

	output, err := app.ReadEntry(author, id)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, output.ID, id)
	assert.Equal(t, output.Text, "hello")
	assert.Equal(t, output.CreatorID, "123")
}

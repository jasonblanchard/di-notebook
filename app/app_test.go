package app

import (
	"testing"

	"github.com/jasonblanchard/di-notebook/store/postgres"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties/assert"
)

func makeApp() (*App, error) {
	reader, err := postgres.NewReader(&postgres.ReaderInput{
		User:     "di",
		Password: "di",
		Dbname:   "di_notebook",
		Host:     "localhost",
	})

	writer, err := postgres.NewWriter(&postgres.WriterInput{
		User:     "di",
		Password: "di",
		Dbname:   "di_notebook",
		Host:     "localhost",
	})

	if err != nil {
		return nil, err
	}

	app := &App{
		Reader: reader,
		Writer: writer,
	}

	return app, nil
}

func TestCreateReadFlow(t *testing.T) {
	app, err := makeApp()

	if err != nil {
		panic(err)
	}

	tester := &Principle{
		Type: "TEST",
	}

	err = app.ResetEntries(tester)

	if err != nil {
		panic(err)
	}

	author := &Principle{
		Type: "USER",
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

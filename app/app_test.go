package app

import (
	"fmt"
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
		Port:     "5432",
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

func createEntries(app *App, creatorID string, n int) error {
	author := &Principal{
		Type: PrincipalUSER,
		ID:   creatorID,
	}

	for i := 0; i < n; i++ {
		_, err := app.StartNewEntry(&StartNewEntryInput{
			Principal: author,
			Text:      fmt.Sprintf("Hello %d", i),
			CreatorID: creatorID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func TestCreateReadFlow(t *testing.T) {
	app, err := makeApp()

	if err != nil {
		panic(err)
	}

	err = app.StoreWriter.DropEntries()

	if err != nil {
		panic(err)
	}

	author := &Principal{
		Type: PrincipalUSER,
		ID:   "123",
	}

	id, err := app.StartNewEntry(&StartNewEntryInput{
		Principal: author,
		Text:      "hello",
		CreatorID: "123",
	})
	if err != nil {
		panic(err)
	}

	output, err := app.ReadEntry(&ReadEntryInput{
		Principal: author,
		ID:        id,
	})

	if err != nil {
		panic(err)
	}

	assert.Equal(t, output.ID, id)
	assert.Equal(t, output.Text, "hello")
	assert.Equal(t, output.CreatorID, "123")
}

func TestUpdateFlow(t *testing.T) {
	app, err := makeApp()

	if err != nil {
		panic(err)
	}

	err = app.StoreWriter.DropEntries()

	if err != nil {
		panic(err)
	}

	author := &Principal{
		Type: PrincipalUSER,
		ID:   "123",
	}

	id, err := app.StartNewEntry(&StartNewEntryInput{
		Principal: author,
		Text:      "hello",
		CreatorID: "123",
	})
	if err != nil {
		panic(err)
	}

	_, err = app.ChangeEntry(&ChangeEntryInput{
		Principal: author,
		ID:        id,
		Text:      "hello updated",
	})

	if err != nil {
		panic(err)
	}

	output, err := app.ReadEntry(&ReadEntryInput{
		Principal: author,
		ID:        id,
	})

	if err != nil {
		panic(err)
	}

	assert.Equal(t, output.ID, id)
	assert.Equal(t, output.Text, "hello updated")
}

func TestListEntries(t *testing.T) {
	app, err := makeApp()
	if err != nil {
		panic(err)
	}

	err = app.StoreWriter.DropEntries()

	if err != nil {
		panic(err)
	}

	err = createEntries(app, "123", 20)

	output, err := app.ListEntries(&ListEntriesInput{
		CreatorID: "123",
		First:     5,
	})

	if err != nil {
		panic(err)
	}

	nums := []int{19, 18, 17, 16, 15}

	assert.Equal(t, len(output), 5)
	for i, o := range output {
		assert.Equal(t, o.Text, fmt.Sprintf("Hello %d", nums[i]))
	}
	// TODO: Check pagination data
	last := output[len(output)-1]
	lastID := last.ID

	// TODO: Get next n after last
	output, err = app.ListEntries(&ListEntriesInput{
		CreatorID: "123",
		First:     5,
		After:     lastID,
	})

	if err != nil {
		panic(err)
	}

	nums = []int{14, 13, 12, 11, 10}

	assert.Equal(t, len(output), 5)
	for i, o := range output {
		assert.Equal(t, o.Text, fmt.Sprintf("Hello %d", nums[i]))
	}

	// TODO: Get n after last where n > what's left
}

package main

import (
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/di_messages/entry"
	"github.com/jasonblanchard/di-notebook/store/postgres"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func makeApp() (*app.App, error) {
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

	app := &app.App{
		StoreReader: reader,
		StoreWriter: writer,
	}

	return app, nil
}

func TestEndToEnd(t *testing.T) {
	a, err := makeApp()
	if err != nil {
		panic(err)
	}

	a.StoreWriter.DropEntries()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	createEntryRequest := &entry.CreateEntryRequest{
		Context: &entry.RequestContext{
			Principal: &entry.Principal{
				Type: entry.Principal_USER,
				Id:   "123",
			},
		},
		Payload: &entry.CreateEntryRequest_Payload{
			Text:      "Hello, world",
			CreatorId: "123",
		},
	}

	data, err := proto.Marshal(createEntryRequest)
	if err != nil {
		panic(err)
	}

	response, err := nc.Request("create.entry", data, time.Second*1)
	if err != nil {
		panic(err)
	}

	createResponseBody := &entry.CreateEntryResponse{}
	err = proto.Unmarshal(response.Data, createResponseBody)
	if err != nil {
		panic(err)
	}

	getEntryRequest := &entry.GetEntryRequest{
		Context: &entry.RequestContext{
			Principal: &entry.Principal{
				Type: entry.Principal_USER,
				Id:   "123",
			},
		},
		Payload: &entry.GetEntryRequest_Payload{
			Id: createResponseBody.Payload.Id,
		},
	}

	getEntryRequestData, err := proto.Marshal(getEntryRequest)
	if err != nil {
		panic(err)
	}

	getResponse, err := nc.Request("get.entry", getEntryRequestData, time.Second*1)

	getResponseBody := &entry.GetEntryResponse{}
	err = proto.Unmarshal(getResponse.Data, getResponseBody)

	assert.Equal(t, getResponseBody.Payload.Text, "Hello, world")
}

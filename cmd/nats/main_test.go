package main

import (
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/di_messages/entry"
	"github.com/jasonblanchard/di-notebook/store/postgres"
	"github.com/nats-io/nats.go"
	"google.golang.org/genproto/googleapis/rpc/code"
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

	readEntryRequest := &notebook.ReadEntryRequest{
		Context: &notebook.RequestContext{
			Principal: &notebook.Principal{
				Type: notebook.Principal_USER,
				Id:   "123",
			},
			TraceId: "666",
		},
		Payload: &notebook.ReadEntryRequest_Payload{
			Id: createResponseBody.Payload.Id,
		},
	}

	readEntryRequestData, err := proto.Marshal(readEntryRequest)
	if err != nil {
		panic(err)
	}

	readEntryResponse, err := nc.Request("notebook.ReadEntry", readEntryRequestData, time.Second*1)
	if err != nil {
		panic(err)
	}

	readResponseBody := &notebook.ReadEntryResponse{}
	err = proto.Unmarshal(readEntryResponse.Data, readResponseBody)
	if err != nil {
		panic(err)
	}

	// assert.Equal(t, int32(code.Code_PERMISSION_DENIED), readResponseBody.GetStatus().GetCode())
	assert.Equal(t, int32(code.Code_OK), readResponseBody.GetStatus().GetCode())
	assert.Equal(t, "Hello, world", readResponseBody.GetPayload().GetText())
	assert.Equal(t, &notebook.NullableTimestamp_Null{}, readResponseBody.GetPayload().GetUpdatedAt().GetValue())
	assert.Equal(t, readResponseBody.GetContext().TraceId, "666")
}

func TestEndToEndError(t *testing.T) {
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

	readEntryRequest := &notebook.ReadEntryRequest{
		Context: &notebook.RequestContext{
			Principal: &notebook.Principal{
				Type: notebook.Principal_USER,
				Id:   "NOT SAME USER",
			},
			TraceId: "666",
		},
		Payload: &notebook.ReadEntryRequest_Payload{
			Id: createResponseBody.Payload.Id,
		},
	}

	readEntryRequestData, err := proto.Marshal(readEntryRequest)
	if err != nil {
		panic(err)
	}

	readEntryResponse, err := nc.Request("notebook.ReadEntry", readEntryRequestData, time.Second*1)
	if err != nil {
		panic(err)
	}

	readResponseBody := &notebook.ReadEntryResponse{}
	err = proto.Unmarshal(readEntryResponse.Data, readResponseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, int32(code.Code_PERMISSION_DENIED), readResponseBody.GetStatus().GetCode())
	assert.Equal(t, "Permission denied", readResponseBody.Status.GetMessage())
}

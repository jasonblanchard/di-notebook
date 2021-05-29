package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/jasonblanchard/di-notebook/pkg/app"
	"github.com/jasonblanchard/di-notebook/pkg/store/postgres"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestReadEntryNotFound(t *testing.T) {
	conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := notebook.NewNotebookClient(conn)

	request := &notebook.ReadEntryGRPCRequest{
		Principal: &notebook.Principal{
			Id:   "1",
			Type: notebook.Principal_USER,
		},
		Payload: &notebook.ReadEntryGRPCRequest_Payload{
			Id: "123",
		},
	}

	ctx := context.TODO()

	_, err = client.ReadEntry(ctx, request)

	status, _ := status.FromError(err)
	assert.Equal(t, status.Code(), codes.NotFound, fmt.Sprintf("%s", status.Message()))
}

func TestCreateAndRead(t *testing.T) {
	a, err := makeApp()
	if err != nil {
		panic(err)
	}

	a.StoreWriter.DropEntries()

	conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := notebook.NewNotebookClient(conn)

	startNewEntryRequest := &notebook.StartNewEntryGRPCRequest{
		Principal: &notebook.Principal{
			Id:   "1",
			Type: notebook.Principal_USER,
		},
		Payload: &notebook.StartNewEntryGRPCRequest_Payload{
			CreatorId: "1",
		},
	}

	ctx := context.TODO()

	startNewEntryResponse, err := client.StartNewEntry(ctx, startNewEntryRequest)

	assert.Nil(t, err)
	assert.NotEmpty(t, startNewEntryResponse.GetPayload().GetId())

	readRequest := &notebook.ReadEntryGRPCRequest{
		Principal: &notebook.Principal{
			Id:   "1",
			Type: notebook.Principal_USER,
		},
		Payload: &notebook.ReadEntryGRPCRequest_Payload{
			Id: startNewEntryResponse.GetPayload().GetId(),
		},
	}

	readResponse, err := client.ReadEntry(ctx, readRequest)
	assert.Nil(t, err)
	assert.Equal(t, startNewEntryResponse.GetPayload().GetId(), readResponse.GetPayload().GetId())
	assert.Equal(t, "", readResponse.GetPayload().GetText())
	assert.Equal(t, "1", readResponse.GetPayload().GetCreatorId())
	assert.NotEmpty(t, readResponse.GetPayload().CreatedAt)
	assert.Nil(t, readResponse.GetPayload().UpdatedAt)
}

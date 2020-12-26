package main

import (
	"context"
	"log"
	"testing"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestReadEntry(t *testing.T) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := notebook.NewNotebookClient(conn)

	request := &notebook.ReadEntryGRPCRequest{
		RequestContext: &notebook.GRPCRequestContext{
			Principal: &notebook.Principal{
				Id:   "1",
				Type: notebook.Principal_USER,
			},
		},
		Payload: &notebook.ReadEntryGRPCRequest_Payload{
			Id: "123",
		},
	}

	ctx := context.TODO()

	_, err = client.ReadEntry(ctx, request)

	status, _ := status.FromError(err)
	assert.Equal(t, status.Code(), codes.NotFound)

	// assert.Nil(t, err)

	// assert.Equal(t, "", response.GetId())
	// assert.Equal(t, "", response.GetText())
}

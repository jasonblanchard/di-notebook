package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestReadEntry(t *testing.T) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := notebook.NewNotebookClient(conn)

	request := &notebook.ReadEntryGRPCRequest{
		Id: "123",
	}

	ctx := context.TODO()

	response, err := client.ReadEntry(ctx, request)

	fmt.Println(response)

	assert.Equal(t, "123", response.GetId())
	assert.Equal(t, "testing, testing", response.GetText())
}

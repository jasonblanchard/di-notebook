package main

import (
	"context"
	"fmt"
	"net"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"google.golang.org/grpc"
)

type notebookServer struct {
	notebook.UnimplementedNotebookServer
}

func (s *notebookServer) ReadEntry(ctx context.Context, request *notebook.ReadEntryGRPCRequest) (*notebook.ReadEntryGRPCResponse, error) {
	response := &notebook.ReadEntryGRPCResponse{
		Id:   "123",
		Text: "testing, testing",
	}
	return response, nil
}

func main() {
	port := "8080"

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	fmt.Println(fmt.Sprintf("Listening on port %s", port))

	grpcServer := grpc.NewServer()
	notebook.RegisterNotebookServer(grpcServer, &notebookServer{})
	grpcServer.Serve(lis)
}

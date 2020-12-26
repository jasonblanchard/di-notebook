package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "", "Config file")
	flag.Parse()

	err := initConfig(cfgFile)
	if err != nil {
		panic(err)
	}

	// TODO: Make configurable
	port := "8080"

	errorHandler := func(p interface{}) (err error) {
		fmt.Println("Oops")
		fmt.Println(err)
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(errorHandler),
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	fmt.Println(fmt.Sprintf("Listening on port %s", port))

	s, err := NewService()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)
	notebook.RegisterNotebookServer(grpcServer, s)
	grpcServer.Serve(lis)
}

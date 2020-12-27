package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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

	defer s.logger.Sync()

	grpc_zap.ReplaceGrpcLoggerV2(s.logger)

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(s.handleError),
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(s.logger),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(s.logger),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	notebook.RegisterNotebookServer(grpcServer, s)
	grpcServer.Serve(lis)
}

package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

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

	s, err := NewService()
	if err != nil {
		panic(err)
	}

	port := s.Port

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	s.Logger.Info(fmt.Sprintf("Listening on port %s 🚀", port))

	defer s.Logger.Sync()

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		s.Logger.Info("/health")
		fmt.Fprintf(w, "ok")
	})

	// TODO: Make this better, configurable and check in healthcheck
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	grpc_zap.ReplaceGrpcLoggerV2(s.Logger)

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(s.handleError),
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(s.Logger),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(s.Logger),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	notebook.RegisterNotebookServer(grpcServer, s)
	grpcServer.Serve(lis)
}

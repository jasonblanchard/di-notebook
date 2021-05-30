package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	notebook "github.com/jasonblanchard/di-apis/gen/pb-go/notebook/v2"
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

	srv, err := NewServer()
	if err != nil {
		panic(err)
	}

	port := srv.Port

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	srv.Logger.Info(fmt.Sprintf("Listening on port %s 🚀", port))

	defer srv.Logger.Sync()

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		srv.Logger.Info("/health")
		fmt.Fprintf(w, "ok")
	})

	// TODO: Make this better, configurable and check in healthcheck
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	grpc_zap.ReplaceGrpcLoggerV2(srv.Logger)

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(srv.handleError),
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(srv.Logger),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(srv.Logger),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	notebook.RegisterNotebookServer(grpcServer, srv)
	grpcServer.Serve(lis)
}

package main

import (
	"flag"
	"fmt"
	"net"

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

	grpcServer := grpc.NewServer()
	notebook.RegisterNotebookServer(grpcServer, s)
	grpcServer.Serve(lis)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	log.Printf("Gin cold start")
	r := echo.New()
	echoLambda = echoadapter.New(r)

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

	r.GET("/api/meta", srv.HandleMeta)
	// r.GET("/api/me", srv.HandleMe)
	// r.GET("/api/v2/entries", srv.HandleListEntries)
	// r.GET("/api/v2/entries/:id", srv.HandleGetEntry)
	// r.POST("/api/v2/entries", srv.HandleCreateEntry)
	// r.PATCH("/api/v2/entries/:id", srv.HandleUpdateEntry)
	// r.DELETE("/api/v2/entries/:id", srv.HandleDeleteEntry)
	// r.POST("/api/v2/entries/:id/undelete", srv.HandleUndeleteEntry)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(fmt.Sprintf("%+v", req))
	return echoLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}

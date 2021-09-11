package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Printf("Gin cold start")
	r := gin.Default()
	ginLambda = ginadapter.New(r)

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
	r.GET("/api/me", srv.HandleMe)
	r.GET("/api/v2/entries", srv.HandleListEntries)
	r.GET("/api/v2/entries/:id", srv.HandleGetEntry)
	r.POST("/api/v2/entries", srv.HandleCreateEntry)
	r.PATCH("/api/v2/entries/:id", srv.HandleUpdateEntry)
	r.DELETE("/api/v2/entries/:id", srv.HandleDeleteEntry)
	r.POST("/api/v2/entries/:id/undelete", srv.HandleUndeleteEntry)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(fmt.Sprintf("%+v", req))
	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}

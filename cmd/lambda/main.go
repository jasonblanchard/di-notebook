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
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(fmt.Sprintf("%+v", req))
	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}

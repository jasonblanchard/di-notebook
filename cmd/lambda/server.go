package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/gin-gonic/gin"
	"github.com/jasonblanchard/di-notebook/pkg/app"
	"github.com/jasonblanchard/di-notebook/pkg/store/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Server struct {
	*app.App
	Logger *zap.Logger
	Port   string
}

func initConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	return nil
}

// NewServer Create a new Server from env
func NewServer() (*Server, error) {
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	database := viper.GetString("DATABASE")
	pretty := viper.GetBool("PRETTY")
	port := viper.GetString("PORT")

	s := &Server{}

	s.Port = port

	db, err := postgres.NewConnection(&postgres.NewConnectionInput{
		User:     dbUser,
		Password: dbPassword,
		Dbname:   database,
		Host:     dbHost,
		Port:     dbPort,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create database")
	}

	reader := &postgres.Reader{
		Db: db,
	}

	writer := &postgres.Writer{
		Db: db,
	}

	s.App = &app.App{
		StoreReader: reader,
		StoreWriter: writer,
	}

	var logger *zap.Logger

	if pretty == true {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create logger")
	}

	s.Logger = logger

	return s, nil
}

func (s *Server) HandleMeta(c *gin.Context) {
	apiGwContext, err := ginLambda.GetAPIGatewayContext(c.Request)
	context := fmt.Sprintf("%+v", apiGwContext)

	requestId := apiGwContext.RequestID
	stage := apiGwContext.Stage

	authorizer := apiGwContext.Authorizer

	version := lambdacontext.FunctionVersion

	authorizationHeader := c.Request.Header["Authorization"]

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, gin.H{
		"context":             context,
		"requestId":           requestId,
		"stage":               stage,
		"authorizer":          authorizer,
		"authorizationHeader": authorizationHeader,
		"version":             version,
	})
}

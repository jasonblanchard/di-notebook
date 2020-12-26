package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/store/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Service service container
type Service struct {
	*app.App
	notebook.UnimplementedNotebookServer
}

func initConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	return nil
}

// NewService Create a new service from env
func NewService() (*Service, error) {
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	database := viper.GetString("DATABASE")

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

	s := &Service{
		App: &app.App{
			StoreReader: reader,
			StoreWriter: writer,
		},
	}

	return s, nil
}

// ReadEntry implements ReadEntry
func (s *Service) ReadEntry(ctx context.Context, request *notebook.ReadEntryGRPCRequest) (*notebook.ReadEntryGRPCResponse, error) {
	id, err := strconv.Atoi(request.GetPayload().GetId())
	if err != nil {
		return nil, MapError(err)
	}

	readEntryInput := &app.ReadEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   request.GetRequestContext().GetPrincipal().Id,
		},
		ID: id,
	}

	entry, err := s.App.ReadEntry(readEntryInput)
	if err != nil {
		return nil, MapError(err)
	}

	response := &notebook.ReadEntryGRPCResponse{
		Payload: &notebook.ReadEntryGRPCResponse_Payload{
			Id:   fmt.Sprintf("%d", entry.ID),
			Text: entry.Text,
		},
	}
	return response, nil
}

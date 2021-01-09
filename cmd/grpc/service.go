package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/store/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Service service container
type Service struct {
	*app.App
	notebook.UnimplementedNotebookServer
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

// NewService Create a new service from env
func NewService() (*Service, error) {
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	database := viper.GetString("DATABASE")
	pretty := viper.GetBool("PRETTY")
	port := viper.GetString("PORT")

	s := &Service{}

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

func (s *Service) handleError(p interface{}) error {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

// ReadEntry implements ReadEntry
func (s *Service) ReadEntry(ctx context.Context, request *notebook.ReadEntryGRPCRequest) (*notebook.ReadEntryGRPCResponse, error) {
	if request.GetPayload().GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is required")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok != true {
		s.Logger.Error("No metadata")
		return nil, status.Error(codes.Unknown, "Error")
	}

	principal, err := getPrincipal(md)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "Error")
	}

	id, err := strconv.Atoi(request.GetPayload().GetId())
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.NotFound, "Not found")
	}

	readEntryInput := &app.ReadEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		ID: id,
	}

	entry, err := s.App.ReadEntry(readEntryInput)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, MapError(err)
	}

	response := &notebook.ReadEntryGRPCResponse{
		Payload: &notebook.ReadEntryGRPCResponse_Payload{
			Id:        fmt.Sprintf("%d", entry.ID),
			CreatorId: entry.CreatorID,
			Text:      entry.Text,
			CreatedAt: timeToProtoTime(entry.CreatedAt),
		},
	}

	if !entry.UpdatedAt.IsZero() {
		response.Payload.UpdatedAt = timeToProtoTime(entry.UpdatedAt)
	}
	return response, nil
}

// StartNewEntry implements StartNewEntry
func (s *Service) StartNewEntry(ctx context.Context, request *notebook.StartNewEntryGRPCRequest) (*notebook.StartNewEntryGRPCResponse, error) {
	if request.GetPayload().GetCreatorId() == "" {
		return nil, status.Error(codes.InvalidArgument, "CreatorId is required")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok != true {
		s.Logger.Error("No metadata")
		return nil, status.Error(codes.Unknown, "Error")
	}

	principal, err := getPrincipal(md)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "Error")

	}

	input := &app.StartNewEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		CreatorID: request.GetPayload().GetCreatorId(),
	}

	id, err := s.App.StartNewEntry(input)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, MapError(err)
	}

	response := &notebook.StartNewEntryGRPCResponse{
		Payload: &notebook.StartNewEntryGRPCResponse_Payload{
			Id: fmt.Sprintf("%d", id),
		},
	}

	return response, nil
}

func timeToProtoTime(time time.Time) *timestamp.Timestamp {
	seconds := time.Unix()

	if time.IsZero() {
		seconds = 0
	}

	return &timestamp.Timestamp{
		Seconds: seconds,
	}
}

func getPrincipal(md metadata.MD) (*notebook.Principal, error) {
	data, ok := md["principal-bin"]
	if ok == false {
		return nil, errors.New("principal key missing from metadata")
	}

	principal := &notebook.Principal{}
	err := proto.Unmarshal([]byte(strings.Join(data, "")), principal)

	if err != nil {
		return nil, errors.New("Error unmarshalling principal")
	}

	return principal, nil
}

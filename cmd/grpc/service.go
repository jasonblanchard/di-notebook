package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	notebook "github.com/jasonblanchard/di-apis/gen/pb-go/notebook/v2"
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/store/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	fmt.Println("Oops")
	fmt.Println(fmt.Sprintf("panic triggered: %v", p))
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

// GetEntry implements GetEntry
func (s *Service) GetEntry(ctx context.Context, request *notebook.GetEntryRequest) (*notebook.Entry, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok != true {
		return nil, status.Error(codes.InvalidArgument, "Missing metadata")
	}

	principal, err := getPrincipal(md)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "Error")
	}

	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is required")
	}

	id, err := strconv.Atoi(request.GetId())
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

	response := &notebook.Entry{
		Id:        fmt.Sprintf("%d", entry.ID),
		CreatorId: entry.CreatorID,
		Text:      entry.Text,
		CreatedAt: timeToProtoTime(entry.CreatedAt),
	}

	if !entry.UpdatedAt.IsZero() {
		response.UpdatedAt = timeToProtoTime(entry.UpdatedAt)
	}
	return response, nil
}

// CreateEntry implements CreateEntry
func (s *Service) CreateEntry(ctx context.Context, request *notebook.CreateEntryRequest) (*notebook.Entry, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok != true {
		return nil, status.Error(codes.InvalidArgument, "Missing metadata")
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
		CreatorID: principal.GetId(),
		Text:      request.GetEntry().GetText(),
	}

	id, err := s.App.StartNewEntry(input)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, MapError(err)
	}

	response := &notebook.Entry{
		Id:   fmt.Sprintf("%d", id),
		Text: request.GetEntry().GetText(),
	}

	return response, nil
}

// ListEntries implements ListEntries
func (s *Service) ListEntries(ctx context.Context, request *notebook.ListEntryRequest) (*notebook.ListEntriesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO")
}

// UpdateEntry implements UpdateEntry
func (s *Service) UpdateEntry(ctx context.Context, request *notebook.UpdateEntryRequest) (*notebook.Entry, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok != true {
		return nil, status.Error(codes.InvalidArgument, "Missing metadata")
	}

	principal, err := getPrincipal(md)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "Error")
	}

	input := &app.ChangeEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		Text: request.GetEntry().GetText(),
	}

	entry, err := s.App.ChangeEntry(input)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, MapError(err)
	}

	response := &notebook.Entry{
		Id:        fmt.Sprintf("%d", entry.ID),
		CreatorId: entry.CreatorID,
		Text:      entry.Text,
		CreatedAt: timeToProtoTime(entry.CreatedAt),
		UpdatedAt: timeToProtoTime(entry.UpdatedAt),
	}

	return response, nil
}

// DeleteEntry implements DeleteEntry
func (s *Service) DeleteEntry(ctx context.Context, request *notebook.DeleteEntryRequest) (*empty.Empty, error) {
	return &empty.Empty{}, status.Error(codes.Unimplemented, "TODO")
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
	bearer, ok := md["authorization"]
	if ok == false {
		return nil, errors.New("principal key missing from metadata")
	}

	id, err := bearerHeaderToID(strings.Join(bearer, ""))

	if err != nil {
		return nil, err
	}

	principal := &notebook.Principal{
		Id: id,
	}

	return principal, nil
}

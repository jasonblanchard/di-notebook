package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	notebook "github.com/jasonblanchard/di-apis/gen/pb-go/notebook/v2"
	"github.com/jasonblanchard/di-notebook/pkg/app"
	"github.com/jasonblanchard/di-notebook/pkg/store/postgres"
	"github.com/nats-io/nats.go"
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
	Logger         *zap.Logger
	Port           string
	NatsConnection *nats.Conn
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
	natsURL := viper.GetString("NATS_URL")
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

	if natsURL != "" {
		nc, err := nats.Connect(natsURL)
		if err != nil {
			return nil, errors.Wrap(err, "NATS initialization failed")
		}
		s.NatsConnection = nc
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
	s.Logger.Error("Captured panic")
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

	getEntryInput := &app.GetEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		ID: id,
	}

	entry, err := s.App.GetEntry(getEntryInput)
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

	input := &app.CreateEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		CreatorID: principal.GetId(),
		Text:      request.GetEntry().GetText(),
	}

	id, err := s.App.CreateEntry(input)
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
	md, ok := metadata.FromIncomingContext(ctx)

	if ok != true {
		return nil, status.Error(codes.InvalidArgument, "Missing metadata")
	}

	principal, err := getPrincipal(md)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "Error")
	}

	var after int
	if request.GetPageToken() == "" {
		after = 0
	} else {
		after, err = strconv.Atoi(request.GetPageToken())
	}

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unknown, "Error")
	}

	input := &app.ListEntriesInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		CreatorID: principal.GetId(),
		First:     int(request.GetPageSize()),
		After:     after, // TODO: Consider a proper opaque token for this, like base64'd id
	}

	output, err := s.App.ListEntries(input)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, MapError(err)
	}

	response := &notebook.ListEntriesResponse{
		NextPageToken: fmt.Sprintf("%d", output.Pagination.EndCursor),
		TotalSize:     int32(output.Pagination.TotalCount),
		HasNextPage:   output.Pagination.HasNextPage,
	}

	for _, entry := range output.Entries {
		responseEntry := &notebook.Entry{
			Id:        fmt.Sprintf("%d", entry.ID),
			Text:      entry.Text,
			CreatorId: entry.CreatorID,
			CreatedAt: timeToProtoTime(entry.CreatedAt),
		}

		if !entry.UpdatedAt.IsZero() {
			responseEntry.UpdatedAt = timeToProtoTime(entry.UpdatedAt)
		}

		response.Entries = append(response.Entries, responseEntry)
	}

	return response, nil
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

	id, err := strconv.Atoi(request.GetId())
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unknown, "Error")
	}

	input := &app.UpdateEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		ID:   id,
		Text: request.GetEntry().GetText(),
	}

	entry, err := s.App.UpdateEntry(input, func(entry *app.Entry) {
		revision, err := EntryToEntryRevision(entry, principal)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		s.NatsConnection.Publish("data.mesh.notebook.v2.EntryRevision", revision)
	})

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
func (s *Service) DeleteEntry(ctx context.Context, request *notebook.DeleteEntryRequest) (*notebook.DeleteEntryResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok != true {
		return nil, status.Error(codes.InvalidArgument, "Missing metadata")
	}

	principal, err := getPrincipal(md)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "Error")
	}

	id, err := strconv.Atoi(request.GetId())
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unknown, "Error")
	}

	input := &app.DeleteEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		ID: id,
	}

	entry, err := s.App.DeleteEntry(input, func(entry *app.Entry) {
		revision, err := EntryToEntryRevision(entry, principal)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		s.NatsConnection.Publish("data.mesh.notebook.v2.EntryRevision", revision)
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, MapError(err)
	}

	response := &notebook.DeleteEntryResponse{
		Entry: &notebook.Entry{
			Id:         fmt.Sprintf("%d", entry.ID),
			CreatorId:  entry.CreatorID,
			Text:       entry.Text,
			CreatedAt:  timeToProtoTime(entry.CreatedAt),
			UpdatedAt:  timeToProtoTime(entry.UpdatedAt),
			DeleteTime: timeToProtoTime(entry.DeleteTime),
		},
	}

	return response, nil
}

// UndeleteEntry implements UndeleteEntry
func (s *Service) UndeleteEntry(ctx context.Context, request *notebook.UndeleteEntryRequest) (*notebook.Entry, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok != true {
		return nil, status.Error(codes.InvalidArgument, "Missing metadata")
	}

	principal, err := getPrincipal(md)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "Error")
	}

	id, err := strconv.Atoi(request.GetId())
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, status.Error(codes.Unknown, "Error")
	}

	input := &app.UndeleteEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   principal.GetId(),
		},
		ID: id,
	}

	entry, err := s.App.UndeleteEntry(input, func(entry *app.Entry) {
		revision, err := EntryToEntryRevision(entry, principal)
		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
		s.NatsConnection.Publish("data.mesh.notebook.v2.EntryRevision", revision)
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, MapError(err)
	}

	response := &notebook.Entry{
		Id:         fmt.Sprintf("%d", entry.ID),
		CreatorId:  entry.CreatorID,
		Text:       entry.Text,
		CreatedAt:  timeToProtoTime(entry.CreatedAt),
		UpdatedAt:  timeToProtoTime(entry.UpdatedAt),
		DeleteTime: timeToProtoTime(entry.DeleteTime),
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

package main

import (
	"os"
	"time"

	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/store/postgres"
	"github.com/jasonblanchard/natsby"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Service service container
type Service struct {
	*app.App
	Logger         *zerolog.Logger
	NATSConnection *nats.Conn
}

// NewService create a new service
// TODO: parameterize input
func NewService() (*Service, error) {
	db, err := postgres.NewConnection(&postgres.NewConnectionInput{
		User:     "di",
		Password: "di",
		Dbname:   "di_notebook",
		Host:     "localhost",
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

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, errors.Wrap(err, "NATS initialization failed")
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.DurationFieldUnit = time.Second

	s := &Service{
		App: &app.App{
			StoreReader: reader,
			StoreWriter: writer,
		},
		NATSConnection: nc,
		Logger:         &logger,
	}

	return s, nil
}

// Run start all listeners
func (s *Service) Run() error {
	engine, err := natsby.New(s.NATSConnection)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize engine")
	}

	engine.Use(natsby.WithLogger(s.Logger))

	engine.Subscribe("create.entry", s.handleCreateEntry)

	engine.Run(func() {
		s.Logger.Info().Msg("Ready to receive messages")
	})

	return nil
}

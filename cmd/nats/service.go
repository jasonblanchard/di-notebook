package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/jasonblanchard/di-notebook/pkg/app"
	"github.com/jasonblanchard/di-notebook/pkg/store/postgres"
	"github.com/jasonblanchard/natsby"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Service service container
type Service struct {
	*app.App
	Logger                                   *zerolog.Logger
	NATSConnection                           *nats.Conn
	FirehoseConnetion                        *firehose.Firehose
	FirehoseEntryRevisionsDeliveryStringName string
}

func initConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	return nil
}

// NewServiceFromEnv create a new service
func NewServiceFromEnv() (*Service, error) {
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	database := viper.GetString("DATABASE")
	natsURL := viper.GetString("NATS_URL")
	debug := viper.GetBool("DEBUG")
	pretty := viper.GetBool("PRETTY")
	awsRegion := viper.GetString("AWS_REGION")
	entryRevisionsDeliveryStreamName := viper.GetString("FIREHOSE_ENTRY_REVISIONS_DELIVERY_STREAM_NAME")

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

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, errors.Wrap(err, "NATS initialization failed")
	}

	logger := initializeLogger(debug, pretty)

	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "AWS session creation failed")
	}
	firehoseConnection := firehose.New(sess, &aws.Config{
		Region: &awsRegion,
	})

	s := &Service{
		App: &app.App{
			StoreReader: reader,
			StoreWriter: writer,
		},
		NATSConnection:                           nc,
		Logger:                                   logger,
		FirehoseConnetion:                        firehoseConnection,
		FirehoseEntryRevisionsDeliveryStringName: entryRevisionsDeliveryStreamName,
	}

	return s, nil
}

func initializeLogger(debug, pretty bool) *zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if pretty == true {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug == true {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.DurationFieldUnit = time.Second
	return &logger
}

// Run start all listeners
func (s *Service) Run() error {
	engine, err := natsby.New(s.NATSConnection)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize engine")
	}

	// TODO: Enable prometheus
	engine.Use(natsby.WithCustomRecovery(errorHandler(s)))
	engine.Use(natsby.WithLogger(s.Logger))

	engine.Subscribe("data.mesh.notebook.v2.EntryRevision", s.handleDebug)

	engine.Run(func() {
		s.Logger.Info().Msg("Ready to receive messages")
	})

	return nil
}

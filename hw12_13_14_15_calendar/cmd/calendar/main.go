package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/app"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/config"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server/http"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Printf("failed to read configuration file: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(config.Logger.Level)

	storage, closeDBConn, err := initStorage(config.Storage.Type, &config.Database)
	if err != nil {
		log.Error("Init storage", "error", err)
		os.Exit(1)
	}

	calendar := app.New(log, storage)

	server := internalhttp.NewServer(log, calendar, &config.Server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if closeDBConn != nil {
			if err := closeDBConn(); err != nil {
				log.Error("Close connection to database", "error", err)
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("Stop http server", "error", err)
		}
	}()

	log.Info("Calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Error("Start http server", "error", err)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

type CloseConnFn func() error

func initStorage(storageType string, config *config.DatabaseConfig) (app.Storage, CloseConnFn, error) {
	switch storageType {
	case storage.InMemory:
		return memorystorage.New(), nil, nil

	case storage.SQL:
		dbConn, err := sqlstorage.NewConnection(config)
		if err != nil {
			return nil, nil, err
		}
		return sqlstorage.New(dbConn), func() error { return dbConn.Close() }, nil

	default:
		return nil, nil, fmt.Errorf("invalid storage type")
	}
}

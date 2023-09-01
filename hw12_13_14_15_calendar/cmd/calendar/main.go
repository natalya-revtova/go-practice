package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/calendar"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/config"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server/http"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage/sql"
	"golang.org/x/exp/slog"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.toml", "Path to configuration file")
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

	calendar := calendar.New(storage)

	serverHTTP := internalhttp.NewServer(log, calendar, &config.ServerHTTP)
	serverGRPC := internalgrpc.NewServer(log, calendar, &config.ServerGRPC)

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

		if err := serverHTTP.Stop(ctx); err != nil {
			log.Error("Stop http server", "error", err)
		}

		serverGRPC.Stop()
	}()

	log.Info("Calendar is running...",
		slog.String("http/server", fmt.Sprintf("%s:%d", config.ServerHTTP.Host, config.ServerHTTP.Port)),
		slog.String("grpc/server", fmt.Sprintf("%s:%d", config.ServerGRPC.Host, config.ServerGRPC.Port)))

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := serverHTTP.Start(); err != nil {
			log.Error("Start http server", "error", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(&config.ServerGRPC); err != nil {
			log.Error("Start grpc server", "error", err)
			cancel()
		}
	}()

	wg.Wait()
}

type CloseConnFn func() error

func initStorage(storageType string, config *config.DatabaseConfig) (calendar.Storage, CloseConnFn, error) {
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

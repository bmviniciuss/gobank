package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bmviniciuss/gobank/internal/api/core/http/personapi"
	"github.com/bmviniciuss/gobank/internal/core/person"
	"github.com/bmviniciuss/gobank/internal/core/person/stores/persondb"
	"github.com/bmviniciuss/gobank/internal/foundation/logger"
	"github.com/bmviniciuss/gobank/internal/foundation/sqldb"
	"github.com/bmviniciuss/gobank/internal/foundation/web"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	var (
		ctx    = context.Background()
		logger = logger.New("person")
	)
	defer logger.Sync()

	if err := run(ctx, logger); err != nil {
		logger.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

type config struct {
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBName     string `env:"DB_NAME" envDefault:"gobank"`
	DBUsername string `env:"DB_USER" envDefault:"gobank"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"gobank"`
	Addr       string `env:"ADDRS" envDefault:":3000"`
}

func run(ctx context.Context, logger *zap.SugaredLogger) error {
	logger.Infof("GOMAXPROCS %d", runtime.GOMAXPROCS(0))
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		logger.With(zap.Error(err)).Error("Got error parsing config")
		return err
	}

	logger.Info("Connecting to database")
	dbpool, err := sqldb.Open(ctx, sqldb.Config{
		User:     cfg.DBUsername,
		Password: cfg.DBPassword,
		Host:     cfg.DBHost,
		Name:     cfg.DBName,
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("Got error ping db connection")
		return err
	}
	logger.Info("Connected to database")

	// -------------------------------------------------------------------------
	// Start API
	logger.Info("Starting Person API server")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	var (
		router        = chi.NewRouter()
		app           = web.NewApp(logger, router)
		personStorer  = persondb.NewStore(logger, dbpool)
		personService = person.NewService(logger, personStorer)
	)

	personapi.Routes(app.Mux, personapi.Config{
		Logger:        logger,
		PersonService: personService,
	})

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: app,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Infof("api router started on host: %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Info("shutdown started")
		defer logger.Infof("shutdown completed sig [%+v]", sig)

		ctx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

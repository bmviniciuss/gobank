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

	"github.com/bmviniciuss/gobank/person/api/http/personapi"
	"github.com/bmviniciuss/gobank/person/core/person"
	"github.com/bmviniciuss/gobank/person/core/person/persondb"
	"github.com/bmviniciuss/gobank/person/foundation/logger"
	"github.com/bmviniciuss/gobank/person/foundation/sqldb"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	var (
		ctx    = context.Background()
		logger = logger.New(logger.Config{Service: "person"})
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
	DBName     string `env:"DB_NAME" envDefault:"person"`
	DBUsername string `env:"DB_USER" envDefault:"person"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"person"`
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
	db, err := sqldb.Open(ctx, sqldb.Config{
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
		personStorer  = persondb.NewStore(logger, db)
		personService = person.NewService(logger, personStorer)
	)

	router.Use(middleware.Logger)
	router.Group(func(r chi.Router) {
		personapi.Routes(r, personapi.Config{
			Logger:        logger,
			PersonService: personService,
		})

	})

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
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

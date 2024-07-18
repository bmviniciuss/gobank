package sqldb

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	Schema       string
	MaxIdleConns int
	MaxOpenConns int
	EnableTLS    bool
}

func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	sslMode := "disable"
	if cfg.EnableTLS {
		sslMode = "require"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	dbpool, err := pgxpool.New(context.Background(), u.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

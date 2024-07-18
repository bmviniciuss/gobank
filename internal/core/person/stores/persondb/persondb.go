package persondb

import (
	"context"
	"errors"
	"fmt"

	"github.com/bmviniciuss/gobank/internal/core/person"
	"github.com/bmviniciuss/gobank/internal/core/person/stores/persondb/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Store struct {
	logger *zap.SugaredLogger
	pool   *pgxpool.Pool
}

func NewStore(
	logger *zap.SugaredLogger,
	pool *pgxpool.Pool,
) *Store {
	return &Store{
		logger: logger,
		pool:   pool,
	}
}

func (s *Store) FindByDocument(ctx context.Context, document string) (*person.Person, error) {
	lggr := s.logger
	db, err := s.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return nil, err
	}
	defer db.Release()
	queries := generated.New(db)
	row, err := queries.FindPersonByDocument(ctx, document)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, person.ErrPersonNotFound
		}
		lggr.With(zap.Error(err)).Error("Got error reading person by document")
		return nil, fmt.Errorf("persondb: %w", err)
	}
	return toPerson(row), nil
}
func (s *Store) Create(ctx context.Context, p *person.Person) error {
	lggr := s.logger
	db, err := s.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return err
	}
	defer db.Release()
	queries := generated.New(db)
	err = queries.InsertPerson(ctx, toInsertPersonParams(p))
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error inserting person")
		return fmt.Errorf("persondb.Create: %w", err)
	}
	return nil
}

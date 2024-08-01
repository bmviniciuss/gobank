package persondb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bmviniciuss/gobank/person/core/person"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Store struct {
	logger *zap.SugaredLogger
	db     *sqlx.DB
}

func NewStore(logger *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		logger: logger,
		db:     db,
	}
}

const (
	findByDocument = `
	SELECT 
		uuid, name, document, active, created_at, updated_at 
	FROM person.person 
	WHERE document = $1 AND active = true
	LIMIT 1;
	`

	insertPerson = `
	INSERT INTO person.person (uuid, name, document, active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6);
	`
)

func (s *Store) FindByDocument(ctx context.Context, document string) (*person.Person, error) {
	lggr := s.logger
	stmt, err := s.db.PreparexContext(ctx, findByDocument)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error preparing statement")
		return nil, err
	}

	var personRow findPersonByDocumentRow
	err = stmt.QueryRowxContext(ctx, document).StructScan(&personRow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			lggr.Error("Person not found")
			return nil, person.ErrPersonNotFound
		}
		lggr.With(zap.Error(err)).Error("Got error querying row")
		return nil, err
	}

	person := personRow.toPerson()
	return &person, nil
}

func (s *Store) Create(ctx context.Context, p *person.Person) error {
	lggr := s.logger
	stmt, err := s.db.PreparexContext(ctx, insertPerson)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error preparing statement")
		return err
	}

	_, err = stmt.ExecContext(ctx, p.ID, p.Name, p.Document, p.Active, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error inserting person")
		return err
	}

	return nil
}

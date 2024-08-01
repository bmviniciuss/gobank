package person

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrPersonNotFound   = errors.New("person not found")
	ErrConflictDocument = errors.New("person with the provided document already exists")
)

type Storer interface {
	FindByDocument(ctx context.Context, document string) (*Person, error)
	FindByID(ctx context.Context, id uuid.UUID) (Person, error)
	Create(ctx context.Context, p *Person) error
}

type Service struct {
	logger *zap.SugaredLogger
	storer Storer
}

func NewService(logger *zap.SugaredLogger, storer Storer) *Service {
	return &Service{
		logger: logger,
		storer: storer,
	}
}

func (s *Service) Create(ctx context.Context, cp CreatePerson) (*Person, error) {
	lggr := s.logger
	lggr.Info("Person creationg started")
	_, err := s.storer.FindByDocument(ctx, cp.Document)
	if err == nil {
		lggr.Error("Person with the provided document already exists")
		return nil, fmt.Errorf("person.Service.Create: %w", ErrConflictDocument)
	}
	if !errors.Is(err, ErrPersonNotFound) {
		lggr.With(zap.Error(err)).Error("Got error reading from store")
		return nil, err
	}

	lggr.Info("Person does not exists, creating one")
	p, err := newPerson(cp.Name, cp.Document)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error creating person entity")
		return nil, err
	}
	err = s.storer.Create(ctx, p)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error storing person entity")
		return nil, err
	}
	return p, nil
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (Person, error) {
	return s.storer.FindByID(ctx, id)
}

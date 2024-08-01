package personapp

import (
	"context"

	"github.com/bmviniciuss/gobank/person/core/person"
	"go.uber.org/zap"
)

type App struct {
	logger        *zap.SugaredLogger
	personService *person.Service
}

func NewApp(logger *zap.SugaredLogger, personService *person.Service) *App {
	return &App{
		logger:        logger,
		personService: personService,
	}
}

func (app *App) Create(ctx context.Context, np NewPerson) (Person, error) {
	lggr := app.logger
	lggr.Info("Creating a new person")
	err := np.Validate()
	if err != nil {
		return Person{}, err
	}
	p, err := app.personService.Create(ctx, person.CreatePerson{
		Name:     np.Name,
		Document: np.Document,
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error creating person")
		return Person{}, err
	}
	var person Person
	person.FromPerson(p)
	return person, nil
}

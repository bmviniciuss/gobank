package personapi

import (
	personapp "github.com/bmviniciuss/gobank/internal/app/core/person"
	"github.com/bmviniciuss/gobank/internal/core/person"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Config struct {
	Logger        *zap.SugaredLogger
	PersonService *person.Service
}

func Routes(app *chi.Mux, cfg Config) {
	api := newAPI(cfg.Logger, personapp.NewApp(cfg.Logger, cfg.PersonService))
	app.Post("/v1/person", api.create)
}

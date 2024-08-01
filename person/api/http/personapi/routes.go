package personapi

import (
	"net/http"
	"time"

	personapp "github.com/bmviniciuss/gobank/person/app/core/person"
	"github.com/bmviniciuss/gobank/person/core/person"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Config struct {
	Logger        *zap.SugaredLogger
	PersonService *person.Service
}

func Routes(r chi.Router, cfg Config) {
	api := newAPI(cfg.Logger, personapp.NewApp(cfg.Logger, cfg.PersonService))
	r.Get("/v1/health",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, map[string]string{
				"status": "ok",
				"time":   time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			})
		}),
	)
	r.Post("/v1/person", api.create)
	r.Get("/v1/person/{id}", api.findByID)
}

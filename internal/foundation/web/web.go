package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.SugaredLogger
	Mux    *chi.Mux
}

func NewApp(
	logger *zap.SugaredLogger,
	mux *chi.Mux,
) *App {
	mux.Use(requestIDMid)
	return &App{
		logger: logger,
		Mux:    mux,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Mux.ServeHTTP(w, r)
}

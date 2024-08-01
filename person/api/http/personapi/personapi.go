package personapi

import (
	"errors"
	"net/http"

	personapp "github.com/bmviniciuss/gobank/person/app/core/person"
	"github.com/bmviniciuss/gobank/person/app/sdk/errs"
	"github.com/bmviniciuss/gobank/person/core/person"
	"github.com/bmviniciuss/gokit/web"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type api struct {
	logger    *zap.SugaredLogger
	personApp *personapp.App
}

func newAPI(logger *zap.SugaredLogger, personApp *personapp.App) *api {
	return &api{
		logger:    logger,
		personApp: personApp,
	}
}

func (api *api) create(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		reqID = uuid.NewString()
		lggr  = api.logger
	)
	var newPerson personapp.NewPerson
	err := web.Decode(r, &newPerson)
	if errs.IsFieldErrors(err) {
		lggr.With(zap.Error(err)).Error("validation error")
		renderError(w, r, web.NewBadRequestErrorResponse(reqID, web.GetFieldErrors(err)))
		return
	}
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding request body")
		renderError(w, r, web.DecodeJSONErrorToResponse(reqID, err))
		return
	}
	p, err := api.personApp.Create(ctx, newPerson)
	if err != nil {
		if errors.Is(err, person.ErrConflictDocument) {
			renderError(w, r, web.NewUnprocessableEntityResponse(reqID, "PERSON_001"))
			return
		}
		lggr.With(zap.Error(err)).Error("Got error creating person")
		return
	}
	renderJSON(w, r, http.StatusCreated, p)
}

func renderError(w http.ResponseWriter, r *http.Request, errRes web.ErrorResponse) {
	renderJSON(w, r, errRes.Status, errRes)
}

func renderJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	render.JSON(w, r, data)
}

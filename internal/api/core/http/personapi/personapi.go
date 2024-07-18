package personapi

import (
	"errors"
	"net/http"

	personapp "github.com/bmviniciuss/gobank/internal/app/core/person"
	"github.com/bmviniciuss/gobank/internal/app/sdk/errs"
	"github.com/bmviniciuss/gobank/internal/core/person"
	"github.com/bmviniciuss/gobank/internal/foundation/web"
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
		reqID = web.GetRequestID(ctx)
		lggr  = api.logger
	)
	var newPerson personapp.NewPerson
	err := web.Decode(r, &newPerson)
	if err != nil {
		if errs.IsFieldErrors(err) {
			lggr.With(zap.Error(err)).Error("validation error")
			ferrs := errs.GetFieldErrors(err)
			fieldErrs := make([]web.FieldError, len(ferrs))
			for i, fe := range ferrs {
				fieldErrs[i] = web.FieldError{
					Field:   fe.Field,
					Message: fe.Message,
				}
			}
			web.RenderError(w, r, web.NewBadRequestErrorResponse(reqID, fieldErrs))
			return
		}
		lggr.With(zap.Error(err)).Error("Got error")
		web.RenderError(w, r, web.DecodeErrorResponse(reqID, err))
		return
	}
	p, err := api.personApp.Create(ctx, newPerson)
	if err != nil {
		if errors.Is(err, person.ErrConflictDocument) {
			web.RenderError(w, r, web.NewUnprocessableEntityResponse(reqID, "PERSON_001"))
			return
		}
		lggr.With(zap.Error(err)).Error("Got error creating person")
		return
	}
	web.Render(w, r, http.StatusCreated, p)
}

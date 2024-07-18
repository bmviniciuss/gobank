package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/render"
)

func Render(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	render.JSON(w, r, data)
}

func RenderError(w http.ResponseWriter, r *http.Request, errorResponse Error) {
	Render(w, r, errorResponse.Status, errorResponse)
}

type Error struct {
	Status int         `json:"-"`
	ID     string      `json:"id"`
	Err    ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// 400
func NewBadRequestErrorResponse(id string, details []FieldError) Error {
	return Error{
		ID:     id,
		Status: http.StatusBadRequest,
		Err: ErrorDetail{
			Code:    "400",
			Message: "Bad Request",
			Details: details,
		},
	}
}

// 404
func NewNotFoundErrorResponse(id string) Error {
	return Error{
		ID:     id,
		Status: http.StatusNotFound,
		Err: ErrorDetail{
			Code:    "404",
			Message: "Not Found",
		},
	}
}

// 404
func NewUnprocessableEntityResponse(id string, code string) Error {
	return Error{
		ID:     id,
		Status: http.StatusUnprocessableEntity,
		Err: ErrorDetail{
			Code:    code,
			Message: "Unprocessable Entity",
		},
	}
}

// 500
func NewInternalServerErrorResponse(id string) Error {
	return Error{
		ID:     id,
		Status: http.StatusInternalServerError,
		Err: ErrorDetail{
			Code:    "500",
			Message: "Internal Server Error",
		},
	}
}

func DecodeErrorResponse(reqID string, err error) Error {
	if err == nil {
		return NewInternalServerErrorResponse(reqID)
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   "body",
				Message: "Invalid JSON",
			},
		})
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   "body",
				Message: "Invalid JSON",
			},
		})
	}
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   unmarshalTypeError.Field,
				Message: "Invalid value",
			},
		})
	}

	if errors.Is(err, io.EOF) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   "body",
				Message: "Empty body",
			},
		})
	}

	return NewInternalServerErrorResponse(reqID)
}

package errs

import (
	"encoding/json"
	"errors"
)

// From: https://github.com/ardanlabs/service/blob/master/app/sdk/errs/errs.go

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// NewFieldsError creates an fields error.
func NewFieldsError(field string, err error) FieldErrors {
	return FieldErrors{
		{
			Field:   field,
			Message: err.Error(),
		},
	}
}

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Encode implements the encoder interface.
func (fe FieldErrors) Encode() ([]byte, string, error) {
	d, err := json.Marshal(fe)
	return d, "application/json", err
}

// Fields returns the fields that failed validation
func (fe FieldErrors) Fields() map[string]string {
	m := make(map[string]string, len(fe))
	for _, fld := range fe {
		m[fld.Field] = fld.Message
	}
	return m
}

// IsFieldErrors checks if an error of type FieldErrors exists.
func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}

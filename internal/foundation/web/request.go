package web

import (
	"fmt"
	"io"
	"net/http"
)

type Decoder interface {
	Decode(data []byte) error
}

type validable interface {
	Validate() error
}

func Decode(r *http.Request, v Decoder) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request: unable to read payload: %w", err)
	}

	if err := v.Decode(data); err != nil {
		return fmt.Errorf("request: decode: %w", err)
	}

	if v, ok := v.(validable); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

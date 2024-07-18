package web

import (
	"net/http"

	"github.com/google/uuid"
)

func requestIDMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		hReqID := r.Header.Get("X-Request-ID")
		if hReqID != "" {
			requestID = hReqID
		}
		ctx := setRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type traceIDKey struct{}

var traceIDKeyValue = traceIDKey{}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), traceIDKeyValue, traceID)
		w.Header().Set("X-Trace-Id", traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

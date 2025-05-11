package web

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func InjectValues(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := &Values{
			TraceID: uuid.NewString(),
			Now:     time.Now(),
		}

		ctx := context.WithValue(r.Context(), Key(), v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

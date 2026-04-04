package middleware

import (
	"context"
	"net/http"
)

func ContextMiddleware(RootCtx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RootCtx.Err() != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithCancel(RootCtx)
		defer cancel()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

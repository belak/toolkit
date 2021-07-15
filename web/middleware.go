package web

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := NewResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				logger.Info().
					Str("method", r.Method).
					Stringer("url", r.URL).
					Int("size", ww.BytesWritten()).
					Int("status", ww.Status()).
					Dur("duration", time.Since(t1)).
					Msg("request complete")
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func RecovererMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					logger.Error().
						Str("stack", string(debug.Stack())).
						Msg("panic while handling request")

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

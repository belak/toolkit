package web

import (
	"net/http"
	"runtime/debug"
	"time"

	"golang.org/x/exp/slog"
)

func LoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := NewResponseWriter(w)

			t1 := time.Now()
			defer func() {
				logger.With(slog.Group("http",
					slog.String("method", r.Method),
					slog.String("url", r.URL.String()),
					slog.Int("size", ww.BytesWritten()),
					slog.Int("status", ww.Status()),
					slog.Duration("duration", time.Since(t1)),
				)).Info("request complete")
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func RecovererMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					logger.With(
						slog.String("stack", string(debug.Stack())),
					).Error("panic while handling request")

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

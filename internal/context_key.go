package internal

import (
	"context"
	"fmt"
	"net/http"
)

type ContextKey string

func (c ContextKey) String() string {
	return fmt.Sprintf("ContextKey<%s>", string(c))
}

// ContextValueMiddleware is a convenience function which stores a static value
// in the request context using the given contextKey.
func ContextValueMiddleware(key ContextKey, val interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, key, val)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

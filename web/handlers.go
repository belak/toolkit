package web

import "net/http"

type ErrorHandler func(ResponseWriter, *http.Request, error)

var DefaultErrorHandler ErrorHandler = defaultErrorHandler

type HandlerFunc func(ResponseWriter, *http.Request) error

func WrapHandler(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(ResponseWriter)
		if !ok {
			rw = NewResponseWriter(w, r.ProtoMajor)
		}
		err := f(rw, r)
		DefaultErrorHandler(rw, r, err)
	}
}

func defaultErrorHandler(w ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	w.Header().Add("Content-Type", "text/plain")

	switch errImpl := err.(type) {
	case StatusError:
		w.WriteHeader(errImpl.Status())
	default:
		w.WriteHeader(500)
	}

	_, _ = w.Write([]byte(err.Error()))
}

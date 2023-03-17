package web

import (
	"net/http"
)

// ResponseWriter is a proxy around an http.ResponseWriter that allows you to
// hook into various parts of the response process.
type ResponseWriter interface {
	http.ResponseWriter

	// Status returns the HTTP status of the request, or 0 if one has not
	// yet been sent.
	Status() int

	// BytesWritten returns the total number of bytes sent to the client.
	BytesWritten() int

	// Unwrap returns the original proxied target.
	Unwrap() http.ResponseWriter
}

// NewResponseWriter wraps an http.ResponseWriter, returning a proxy that allows
// you to hook into various parts of the response process. It must be used with
// http.ResponseController if you need access to additional interfaces.
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &wrappedWriter{ResponseWriter: w}
}

// wrappedWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface, while providing an Unwrap method returning the
// original http.ResponseWriter. This allows it to work properly with
// http.ResponseController without having to support a bunch of different
// permutations of interfaces.
type wrappedWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
}

func (b *wrappedWriter) Status() int                 { return b.code }
func (b *wrappedWriter) BytesWritten() int           { return b.bytes }
func (b *wrappedWriter) Unwrap() http.ResponseWriter { return b.ResponseWriter }

func (b *wrappedWriter) WriteHeader(code int) {
	if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
	}
	b.ResponseWriter.WriteHeader(code)
}

func (b *wrappedWriter) Write(buf []byte) (int, error) {
	b.maybeWriteHeader()
	n, err := b.ResponseWriter.Write(buf)
	b.bytes += n
	return n, err
}

func (b *wrappedWriter) maybeWriteHeader() {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}
}

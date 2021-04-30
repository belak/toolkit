package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeStatus(w http.ResponseWriter, status int, contentType string) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(status)
}

func Stringf(w http.ResponseWriter, format string, args ...interface{}) error {
	writeStatus(w, 200, "")
	_, _ = fmt.Fprintf(w, format, args...)
	return nil
}

func String(w http.ResponseWriter, s string) error {
	writeStatus(w, 200, "")
	_, _ = w.Write([]byte(s))
	return nil
}

func JSON(w http.ResponseWriter, v interface{}) error {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	writeStatus(w, 200, "application/json; charset=utf-8")
	_, _ = w.Write(j)
	return nil
}

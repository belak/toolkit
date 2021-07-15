package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteStatus(w http.ResponseWriter, status int, contentType string) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(status)
}

func RenderStringf(w http.ResponseWriter, format string, args ...interface{}) error {
	WriteStatus(w, 200, "")
	_, _ = fmt.Fprintf(w, format, args...)
	return nil
}

func RenderString(w http.ResponseWriter, s string) error {
	WriteStatus(w, 200, "")
	_, _ = w.Write([]byte(s))
	return nil
}

func RenderJSON(w http.ResponseWriter, v interface{}) error {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	WriteStatus(w, 200, "application/json; charset=utf-8")
	_, _ = w.Write(j)
	return nil
}

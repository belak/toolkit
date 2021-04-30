package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/monoculum/formam"
)

// Creating a decoder which uses the json tag rather than the default `formam`
// allows us to share declarations with the encoding/json package.
var defaultDecoder = formam.NewDecoder(&formam.DecoderOptions{TagName: "json"})

func Decode(r *http.Request, target interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if i := strings.Index(contentType, ";"); i >= 0 {
		contentType = contentType[:i]
	}

	var err error

	if r.Method == http.MethodGet {
		err = defaultDecoder.Decode(r.URL.Query(), target)
	} else if contentType == "application/json" {
		err = json.NewDecoder(r.Body).Decode(target)
	} else if contentType == "application/x-www-form-urlencoded" {
		err = r.ParseForm()
		if err == nil {
			err = defaultDecoder.Decode(r.Form, target)
		}
	} else if contentType == "multipart/form-data" {
		err = r.ParseMultipartForm(32 << 20) // This comes from http.defaultMaxMemory
		if err == nil {
			err = defaultDecoder.Decode(r.Form, target)
		}
	} else {
		return NewError(http.StatusUnsupportedMediaType, fmt.Errorf("unable to handle Content-Type %q", contentType))
	}

	if err != nil {
		return NewError(http.StatusBadRequest, err)
	}

	return nil
}

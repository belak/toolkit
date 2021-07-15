package web

import (
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/rs/zerolog"
)

type TemplateSet struct {
	root *template.Template
}

func NewTemplateSet(rootFS fs.FS, logger zerolog.Logger) (*TemplateSet, error) {
	// Because we want to recursively load the templates, we can't use
	// template.ParseFS, so we walk the dir manually.
	t := template.New("")
	err := fs.WalkDir(rootFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".html") && !strings.HasSuffix(d.Name(), ".htm") {
			return nil
		}

		logger.Debug().Msgf("Loading template from %q", path)

		data, err := fs.ReadFile(rootFS, path)
		if err != nil {
			return err
		}

		t, err = t.New(path).Parse(string(data))
		if err != nil {
			return err
		}

		return nil
	})

	return &TemplateSet{
		root: t,
	}, err
}

func (ts *TemplateSet) Render(w io.Writer, name string, v interface{}) error {
	return ts.root.ExecuteTemplate(w, name, v)
}

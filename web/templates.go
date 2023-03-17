package web

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/belak/toolkit/internal"
	"github.com/belak/toolkit/log"
	"golang.org/x/exp/slog"
)

const templateSetContextKey internal.ContextKey = "TemplateSet"

type TemplateSet struct {
	templates map[string]*template.Template
}

// Because the fs.WalkDir call is very similar between loading partials/layouts
// and loading page templates, this function has a few flags which allow it to
// react in different ways, the most important of which is updateBase, which
// updates the baseTemplate rather than cloning it.
func loadTemplates(logger *slog.Logger, rootFS fs.FS, dir string, baseTemplate *template.Template, updateBase bool) (map[string]*template.Template, error) {
	ret := make(map[string]*template.Template)

	err := fs.WalkDir(rootFS, dir, func(target string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ensure we only operate on files ending in .html
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}

		t := baseTemplate

		if !updateBase {
			// Clone the baseTemplate so the page template can have access to all
			// layouts and includes.
			t, err = baseTemplate.Clone()
			if err != nil {
				return err
			}
		}

		// Because this could technically be run on Windows, we need to sanitize
		// the path name so it's consistent between platforms.
		templateName := filepath.Base(target)

		// Skip any templates which have already been loaded. This makes it
		// possible to have a template directory layout where we preload
		// includes and layouts, but aren't required to put everything else in a
		// separate folder.
		if baseTemplate.Lookup(templateName) != nil {
			return nil
		}

		logger.Debug("Loading template %s from %q", templateName, target)

		// Here's another footgun - when using ParseFS, the template package
		// lops off the directory name and just uses the basename of the file.
		// Because of that, in order to get template names with the directory in
		// the name, we have to read the data and call parse ourselves.
		data, err := fs.ReadFile(rootFS, target)
		if err != nil {
			return err
		}

		t, err = t.New(templateName).Parse(string(data))
		if err != nil {
			return err
		}

		ret[templateName] = t

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func NewTemplateSet(logger *slog.Logger, rootFS fs.FS, funcs ...template.FuncMap) (*TemplateSet, error) {
	ts := &TemplateSet{
		templates: make(map[string]*template.Template),
	}

	baseTemplate := template.New("")

	// Add any common template funcs we care about - we currently add in all the
	// sprig helpers . Note that functions need to be set up before loading
	// templates, or loading the templates will error.
	baseTemplate.Funcs(sprig.FuncMap())

	// Set up any additional built-in functions.
	baseTemplate.Funcs(template.FuncMap{
		"hasField": templateHasField,
	})

	for _, fMap := range funcs {
		baseTemplate.Funcs(fMap)
	}

	// Because the stdlib html/template package has a number of issues, we need
	// to do the parsing in multiple passes, walking the tree every time. First,
	// we need to build a base template which has all the layouts, includes, and
	// functions. Next, we need to load each page template as a separate
	// template so inheritance works how we'd expect.

	_, err := loadTemplates(logger, rootFS, "includes", baseTemplate, true)
	if err != nil {
		return nil, err
	}

	_, err = loadTemplates(logger, rootFS, "layouts", baseTemplate, true)
	if err != nil {
		return nil, err
	}

	// Walk the pages directory and attempt to parse any templates as top-level
	// templates.
	ts.templates, err = loadTemplates(logger, rootFS, ".", baseTemplate, false)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (ts *TemplateSet) Execute(w io.Writer, name string, data interface{}) error {
	t, ok := ts.templates[name]
	if !ok {
		return fmt.Errorf("unknown page template %q", name)
	}

	return t.ExecuteTemplate(w, name, data)
}

func TemplateMiddleware(ts *TemplateSet) func(http.Handler) http.Handler {
	return internal.ContextValueMiddleware(templateSetContextKey, ts)
}

func ExtractTemplates(ctx context.Context) *TemplateSet {
	if ts, ok := ctx.Value(templateSetContextKey).(*TemplateSet); ok {
		return ts
	}

	panic("no template set in context")
}

func Render(ctx context.Context, w io.Writer, name string, data interface{}) {
	logger := log.ExtractLogger(ctx).With(slog.String("template_name", name))
	logger.Debug("rendering template")

	templates := ExtractTemplates(ctx)
	err := templates.Execute(w, name, data)
	if err != nil {
		logger.Error("failed to render template", err)
	}
}

func templateHasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

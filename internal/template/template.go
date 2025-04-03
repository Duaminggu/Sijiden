package template

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// Interface agar bisa dipakai Echo
type TemplateRenderer struct {
	devMode  bool
	compiled *template.Template
}

// Public method untuk buat renderer
func SijidenRenderer(devMode bool) *TemplateRenderer {
	r := &TemplateRenderer{devMode: devMode}
	if !devMode {
		r.compiled = compileTemplates()
	}
	return r
}

// Method untuk Render
func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	var tmpl *template.Template

	if t.devMode {
		// Always reload in dev mode
		tmpl = compileTemplates()
	} else {
		tmpl = t.compiled
	}

	return tmpl.ExecuteTemplate(w, name, data)
}

// Utility function: load & parse all templates
func compileTemplates() *template.Template {
	tmpl := template.New("")
	viewsDir := "views"

	filepath.WalkDir(viewsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			relPath, _ := filepath.Rel(viewsDir, path)
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = tmpl.New(filepath.ToSlash(relPath)).Parse(string(content))
			return err
		}
		return nil
	})

	return tmpl
}

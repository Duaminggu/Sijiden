package template

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func SijidenRenderer() *Template {
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

	return &Template{
		templates: tmpl,
	}
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

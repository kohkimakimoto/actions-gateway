package renderer

import (
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/server/csrf"
	"github.com/kohkimakimoto/actions-gateway/version"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"io/fs"
)

type Renderer struct {
	templates *template.Template
}

func New(fs fs.FS, patterns ...string) *Renderer {
	r := &Renderer{}
	r.templates = template.Must(template.New("T").ParseFS(fs, patterns...))
	return r
}

func (r *Renderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	viewContext, isMap := data.(map[string]any)
	if !isMap {
		return fmt.Errorf("data must be a map[string]any")
	}

	// Add global variables
	viewContext["hash"] = version.ShortCommitHash
	if token := csrf.GetToken(c); token != "" {
		viewContext["csrf"] = template.HTML(fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`, csrf.FormInputName, token))
	}

	return r.templates.ExecuteTemplate(w, name, viewContext)
}

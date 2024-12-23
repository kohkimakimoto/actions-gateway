package renderer

import (
	"bytes"
	"embed"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed testdata/views
var viewsFS embed.FS

func TestRenderer_Render(t *testing.T) {
	r := New(viewsFS, "testdata/views/*.html")
	var buf bytes.Buffer

	e := echo.New()
	ctx := e.NewContext(nil, nil)

	err := r.Render(&buf, "test.html", map[string]any{
		"Name": "World",
	}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!\n", buf.String())
}

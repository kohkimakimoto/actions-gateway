package handlers

import (
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/router"
	openapidocs "github.com/kohkimakimoto/echo-openapidocs"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DocsHandler(r *router.Router) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := r.GetActiveSession(auth.MustGetClient(c))
		if sess == nil {
			return c.String(http.StatusServiceUnavailable, "Your client is not connected to the server")
		}

		return openapidocs.ScalarDocumentsHandler(openapidocs.ScalarConfig{
			Spec: sess.Spec(),
		})(c)
	}
}

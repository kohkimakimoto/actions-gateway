package handlers

import (
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/router"
	"github.com/kohkimakimoto/actions-gateway/server/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDocsHandler(t *testing.T) {
	t.Run("returns status service unavailable when a corresponding session is not found", func(t *testing.T) {
		e := testutil.NewEchoInstance(t)
		r := router.New()
		e.GET("/docs", DocsHandler(r), func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// set a client object to the context for testing
				auth.SetClient(c, &auth.Client{
					Id: "00000000-0000-0000-0000-000000000001",
				})
				return next(c)
			}
		})

		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, "Your client is not connected to the server", rec.Body.String())
	})

	// TODO: Add more tests
}

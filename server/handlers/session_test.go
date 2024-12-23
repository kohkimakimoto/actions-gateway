package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/config"
	"github.com/kohkimakimoto/actions-gateway/server/router"
	"github.com/kohkimakimoto/actions-gateway/server/testutil"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionNewHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := testutil.NewEchoInstance(t)
		cfg := &config.Config{
			URL: "http://localhost:18800",
		}
		r := router.New(router.WithSessionId("00000000-0000-0000-0000-000000000002"))
		ct := &auth.Client{
			Id: "00000000-0000-0000-0000-000000000001",
		}
		e.POST("/api/session/new", SessionNewHandler(cfg, r), func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// set a client object to the context for testing
				auth.SetClient(c, ct)
				return next(c)
			}
		})

		reqBody := &types.SessionNewRequest{
			Actions: []string{"action1", "action2"},
			Spec:    "",
		}
		b, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/session/new", bytes.NewBuffer(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		resp := &types.SessionNewResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), resp)
		assert.NoError(t, err)
		assert.Equal(t, "ws://localhost:18800/api/session/connect/00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002", resp.URL)
	})

	t.Run("invalid request", func(t *testing.T) {
		e := testutil.NewEchoInstance(t)
		e.HTTPErrorHandler = HTTPErrorHandler
		cfg := &config.Config{}
		r := router.New()

		e.POST("/api/session/new", SessionNewHandler(cfg, r))

		req := httptest.NewRequest(http.MethodPost, "/api/session/new", bytes.NewBufferString(`not json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		resp := &types.ErrorResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), resp)
		assert.NoError(t, err)
		assert.Equal(t, "Syntax error: offset=2, error=invalid character 'o' in literal null (expecting 'u')", resp.Error)
		// t.Logf("response: %s", body)
	})

	t.Run("session already exists", func(t *testing.T) {
		e := testutil.NewEchoInstance(t)
		e.HTTPErrorHandler = HTTPErrorHandler
		cfg := &config.Config{
			URL: "http://localhost:18800",
		}
		r := router.New(router.WithSessionId("00000000-0000-0000-0000-000000000002"))
		ct := &auth.Client{
			Id: "00000000-0000-0000-0000-000000000001",
		}
		// create a new session before calling the handler
		_, _ = r.NewSession(ct, []string{"action1", "action2"}, "")

		e.POST("/api/session/new", SessionNewHandler(cfg, r), func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// set a client object to the context for testing
				auth.SetClient(c, ct)
				return next(c)
			}
		})

		reqBody := &types.SessionNewRequest{
			Actions: []string{"action1", "action2"},
			Spec:    "",
		}
		b, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/session/new", bytes.NewBuffer(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		resp := &types.ErrorResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), resp)
		assert.NoError(t, err)
		assert.Equal(t, "Session already exists", resp.Error)
	})
}

// TODO: websocket connection test

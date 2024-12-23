package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testToken:
//
//	{
//	  "alg": "HS256",
//	  "typ": "JWT"
//	}
//
//	{
//	  "iat": 1728912936,
//	  "sub": "01928b3d-cebd-79fa-bd37-9701a24dabdf"
//	}
const testToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3Mjg5MTI5MzYsInN1YiI6IjAxOTI4YjNkLWNlYmQtNzlmYS1iZDM3LTk3MDFhMjRkYWJkZiJ9.twfU54hw8vKqwUAooJqSP3xVBD3OYxmaX6GKY37v8aA`
const testSecret = `12345678901234567890123456789012`

func TestMiddlewareWithConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		k, err := LoadKeyString(testSecret)
		assert.NoError(t, err)

		h := MiddlewareWithConfig(MiddlewareConfig{
			Key: k,
		})(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err = h(c)
		assert.NoError(t, err)

		ct := MustGetClient(c)
		assert.NotNil(t, ct)
		assert.Equal(t, "01928b3d-cebd-79fa-bd37-9701a24dabdf", ct.Id)
	})

	t.Run("error unauthorized: invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		k, err := LoadKeyString("another_key_12345678901234567890")
		assert.NoError(t, err)

		h := MiddlewareWithConfig(MiddlewareConfig{
			Key: k,
		})(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err = h(c)
		assert.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})

	t.Run("error unauthorized: no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		h := MiddlewareWithConfig(MiddlewareConfig{
			Key: nil,
		})(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)
		assert.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})
}

func TestGetClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(nil, nil)
		c.Set(clientKey, &Client{
			Id: "dummy-test-id",
		})
		ct, err := GetClient(c)
		assert.NoError(t, err)
		assert.NotNil(t, ct)
		assert.Equal(t, "dummy-test-id", ct.Id)
	})

	t.Run("error", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(nil, nil)
		ct, err := GetClient(c)
		assert.Error(t, err)
		assert.Nil(t, ct)
		assert.Equal(t, ErrClientNotFound, err)
	})
}

func TestMustGetClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(nil, nil)
		c.Set(clientKey, &Client{
			Id: "dummy-test-id",
		})
		ct := MustGetClient(c)
		assert.NotNil(t, ct)
		assert.Equal(t, "dummy-test-id", ct.Id)
	})

	t.Run("panic", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(nil, nil)
		assert.Panics(t, func() {
			MustGetClient(c)
		})
	})
}

func TestBasicAuthMiddleware(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth(testToken, "")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		k, err := LoadKeyString(testSecret)
		assert.NoError(t, err)

		h := BasicAuthMiddleware(BasicAuthMiddlewareConfig{
			Key: k,
		})(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err = h(c)
		assert.NoError(t, err)

		ct := MustGetClient(c)
		assert.NotNil(t, ct)
		assert.Equal(t, "01928b3d-cebd-79fa-bd37-9701a24dabdf", ct.Id)
	})

	t.Run("error unauthorized: invalid username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth(testToken, "")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		k, err := LoadKeyString("another_key_12345678901234567890")
		assert.NoError(t, err)

		h := BasicAuthMiddleware(BasicAuthMiddlewareConfig{
			Key: k,
		})(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err = h(c)
		assert.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})
}

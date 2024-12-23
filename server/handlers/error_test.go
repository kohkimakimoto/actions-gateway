package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/kohkimakimoto/actions-gateway/server/testutil"
	mock_logger "github.com/kohkimakimoto/actions-gateway/server/testutil/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpErrorHandler(t *testing.T) {
	t.Run("server error", func(t *testing.T) {
		e := testutil.NewEchoInstance(t)
		e.HTTPErrorHandler = HTTPErrorHandler

		returnedError := errors.New("server error")

		// testing logging function is called with the returned error.
		mLogger := mock_logger.NewMockLogger(gomock.NewController(t))
		mLogger.EXPECT().Errorf("%+v", returnedError).Times(1)
		e.Logger = mLogger

		e.GET("/", func(c echo.Context) error {
			return returnedError
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		// The error message is general http error message that doesn't include the original error message.
		assert.JSONEq(t, `{"error":"Internal Server Error"}`, rec.Body.String())
	})

	t.Run("server error with custom error message", func(t *testing.T) {
		e := testutil.NewEchoInstance(t)
		e.HTTPErrorHandler = HTTPErrorHandler

		// Use echo.NewHTTPError to return a custom error message.
		returnedError := echo.NewHTTPError(http.StatusInternalServerError, "Unexpected error")

		// testing logging function is called with the returned error.
		mLogger := mock_logger.NewMockLogger(gomock.NewController(t))
		mLogger.EXPECT().Errorf("%+v", returnedError).Times(1)
		e.Logger = mLogger

		e.GET("/", func(c echo.Context) error {
			return returnedError
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		// The error message is the custom error message.
		assert.JSONEq(t, `{"error":"Unexpected error"}`, rec.Body.String())
	})
}

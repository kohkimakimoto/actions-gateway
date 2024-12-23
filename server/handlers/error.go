package handlers

import (
	"errors"
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	// try to convert error to HTTP error
	var he *echo.HTTPError
	if !errors.As(err, &he) {
		// not HTTP error, create internal server error
		he = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: err,
		}
	}

	if he.Code >= 500 {
		// log 500 errors
		c.Logger().Errorf("%+v", err)
	}

	if c.Request().Method == echo.HEAD {
		// see https://github.com/labstack/echo/issues/608
		_ = c.NoContent(he.Code)
		return
	}

	_ = c.JSON(he.Code, &types.ErrorResponse{
		Error: fmt.Sprintf("%v", he.Message),
	})
}

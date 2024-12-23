package handlers

import (
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/version"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RootHandler(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("Actions Gateway is running. version: %s (hash: %s)", version.Version, version.CommitHash))
}

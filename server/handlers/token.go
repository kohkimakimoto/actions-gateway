package handlers

import (
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/config"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

func NewTokenHandler(tokenGenerator *auth.TokenGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := tokenGenerator.NewTokenAsJWTString()
		if err != nil {
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}

		return c.JSON(http.StatusOK, &types.NewTokenResponse{
			Token: token,
		})
	}
}

func NewTokenPageHandler(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "new_token.html", map[string]any{
			"baseURL": cfg.URL,
		})
	}
}

func NewTokenCreateHandler(cfg *config.Config, tokenGenerator *auth.TokenGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := tokenGenerator.NewTokenAsJWTString()
		if err != nil {
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}

		return c.Render(http.StatusOK, "new_token.html", map[string]any{
			"baseURL": cfg.URL,
			"token":   token,
		})
	}
}

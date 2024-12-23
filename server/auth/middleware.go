package auth

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Client represents a client that is authenticated by a JWT token.
type Client struct {
	Id string
}

// clientKey is the key used to store the client in the Echo context.
const clientKey = "client"

// MiddlewareConfig is the configuration for MiddlewareWithConfig.
type MiddlewareConfig struct {
	Key jwk.Key
}

// MiddlewareWithConfig is an Echo middleware that handles JWT authentication.
func MiddlewareWithConfig(config MiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			// Extract the JWT token from the Authorization header
			tokenString := authHeader[len("Bearer "):]
			// Parse and validate the token using the secret key
			token, err := jwt.ParseString(tokenString, jwt.WithKey(jwa.HS256, config.Key))
			if err != nil {
				// If the token is invalid or the signature does not match, return Unauthorized
				return echo.ErrUnauthorized
			}

			// Set the client in the context
			SetClient(c, &Client{
				Id: token.Subject(),
			})
			return next(c)
		}
	}
}

func SetClient(c echo.Context, client *Client) {
	c.Set(clientKey, client)
}

var ErrClientNotFound = errors.New("client not found in the Echo context")

func GetClient(c echo.Context) (*Client, error) {
	client, ok := c.Get(clientKey).(*Client)
	if !ok {
		return nil, ErrClientNotFound
	}
	return client, nil
}

func MustGetClient(c echo.Context) *Client {
	client, err := GetClient(c)
	if err != nil {
		panic(err)
	}
	return client
}

type BasicAuthMiddlewareConfig struct {
	Key jwk.Key
}

func BasicAuthMiddleware(config BasicAuthMiddlewareConfig) echo.MiddlewareFunc {
	validator := func(username, password string, c echo.Context) (bool, error) {
		// Use the username as the JWT token. The password is not used.
		token, err := jwt.ParseString(username, jwt.WithKey(jwa.HS256, config.Key))
		if err != nil {
			return false, nil
		}

		// Set the client in the context
		SetClient(c, &Client{
			Id: token.Subject(),
		})
		return true, nil
	}

	return middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: validator,
	})
}

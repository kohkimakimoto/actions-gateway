package csrf

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	FormInputName = "_csrf_token"
	HeaderName    = "X-CSRF-Token"
	ContextKey    = "csrf"
	TokenLookup   = fmt.Sprintf("form:%s,header:%s", FormInputName, HeaderName)
)

// GetToken returns the CSRF token from the context.
func GetToken(c echo.Context) string {
	token, ok := c.Get(ContextKey).(string)
	if !ok {
		return ""
	}
	return token
}

func Middleware() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: TokenLookup,
		ContextKey:  ContextKey,
	})
}

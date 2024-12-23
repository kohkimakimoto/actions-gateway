package handlers

import (
	"context"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/router"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func FetchActionHandler(r *router.Router, aFactory *router.ActionMessageFactory) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		sess := r.GetActiveSession(auth.MustGetClient(c))
		if sess == nil {
			return c.String(http.StatusServiceUnavailable, "The session is not active")
		}

		if !sess.IsActionExist(name) {
			return c.String(http.StatusNotFound, "The action is not found")
		}

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}

		// create a new action message
		msg, err := aFactory.NewMessage(name, string(body))
		if err != nil {
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}
		// make sure to free the result channel
		defer sess.FreeResultChannel(msg.Id)

		// Allocate a result channel to receive the action result
		resultChan := sess.AllocateResultChannel(msg.Id)

		// Get a websocket connection from the session
		conn := sess.Conn()
		if conn == nil {
			// Internal server error. The stack trace should be captured.
			return errors.New("the session is active but the websocket connection is nil")
		}

		// Send the action message to the client
		if err := conn.WriteJSON(msg); err != nil {
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second) // Set your desired timeout duration
		defer cancel()

		// Wait for the action result or timeout
		select {
		case result := <-resultChan:
			if result.Status == types.ActionResultStatusSuccess {
				if strings.HasPrefix(result.Body, "{") {
					return c.JSONBlob(http.StatusOK, []byte(result.Body))
				} else {
					return c.String(http.StatusOK, result.Body)
				}
			} else {
				// The action message was sent successfully, but its execution failed.
				// The system returns an internal server error but does not log the error
				// because it is not the fault of the Action Gateway Server.
				// Therefore, the handler returns nil.
				if result.Body != "" {
					return c.String(http.StatusInternalServerError, result.Body)
				}
				return c.String(http.StatusInternalServerError, "The action execution failed")
			}
		case <-ctx.Done():
			// The action execution timed out.
			return c.String(http.StatusInternalServerError, "The action execution timeout")
		}
	}
}

func NotifyActionResultHandler(r *router.Router) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := r.GetActiveSession(auth.MustGetClient(c))
		if sess == nil {
			return c.String(http.StatusServiceUnavailable, "The session is not active")
		}

		result := &types.ActionResult{}
		if err := c.Bind(result); err != nil {
			return err
		}

		if err := sess.HandleActionResult(result); err != nil {
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

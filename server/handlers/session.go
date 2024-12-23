package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/config"
	"github.com/kohkimakimoto/actions-gateway/server/router"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func SessionNewHandler(cfg *config.Config, r *router.Router) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &types.SessionNewRequest{
			Actions: []string{},
		}
		if err := c.Bind(req); err != nil {
			return err
		}

		// create a new session
		client := auth.MustGetClient(c)
		sess, err := r.NewSession(client, req.Actions, req.Spec)
		if err != nil {
			var sessErr *router.SessionError
			if errors.As(err, &sessErr) {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, sessErr.Error()).SetInternal(sessErr)
			}
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}

		c.Logger().Infof("New session created: %s", sess.Key())

		return c.JSON(http.StatusOK, &types.SessionNewResponse{
			URL: cfg.WebSocketURL() + sess.ConnectPath(),
		})
	}
}

const (
	pingPeriod = 10 * time.Second
	pongWait   = 15 * time.Second
)

func SessionConnectHandler(r *router.Router) echo.HandlerFunc {
	return func(c echo.Context) error {
		cId := c.Param("client_id")
		sId := c.Param("session_id")

		// check the client
		client := auth.MustGetClient(c)
		if client.Id != cId {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Client is invalid")
		}

		sess, err := r.ActivateSession(c.Response(), c.Request(), client, sId)
		if err != nil {
			var sessErr *router.SessionError
			if errors.As(err, &sessErr) {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, sessErr.Error()).SetInternal(sessErr)
			}
			// Internal server error. The stack trace should be captured.
			return errors.WithStack(err)
		}
		defer func() {
			r.CloseSession(sess)
			c.Logger().Infof("Session closed: %s", sess.Key())
		}()

		c.Logger().Infof("Session activated: %s", sess.Key())

		// websocket connection
		conn := sess.Conn()

		conn.SetPongHandler(func(appData string) error {
			rAddr := conn.RemoteAddr()
			c.Logger().Debugf("Pong message received from %s", rAddr)
			return conn.SetReadDeadline(time.Now().Add(pongWait))
		})

		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					c.Logger().Infof("Failed to send Ping message: %v", err)
					return
				}
			}
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// terminate the session
				c.Logger().Infof("Session disconnected: %s, %v", sId, err)
				break
			}
		}

		return nil
	}

}

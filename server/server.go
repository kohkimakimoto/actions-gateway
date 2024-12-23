package server

import (
	"context"
	"embed"
	"errors"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/config"
	"github.com/kohkimakimoto/actions-gateway/server/csrf"
	"github.com/kohkimakimoto/actions-gateway/server/handlers"
	"github.com/kohkimakimoto/actions-gateway/server/renderer"
	"github.com/kohkimakimoto/actions-gateway/server/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

var (
	//go:embed public
	publicFS embed.FS
	//go:embed resources/views
	viewsFS embed.FS
)

func Start(cfg *config.Config) error {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Server.Addr = cfg.Addr
	e.Debug = cfg.Debug
	// renderer
	e.Renderer = renderer.New(viewsFS, "resources/views/*.html")

	if e.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}
	e.Logger.SetHeader(`{"time":"${time_rfc3339}","type":"app","level":"${level}"}`)

	e.HTTPErrorHandler = handlers.HTTPErrorHandler

	// ----------------------------------------------------------------
	// Global objects
	// ----------------------------------------------------------------

	// router
	r := router.New()
	// key
	key, err := auth.LoadKeyString(cfg.Secret)
	if err != nil {
		return err
	}
	// token generator
	tokenGenerator := auth.NewTokenGenerator(key)
	// action message factory
	aFactory := router.NewActionMessageFactory()

	// ----------------------------------------------------------------
	// middleware
	// ----------------------------------------------------------------
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogLevel: log.ERROR,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","type":"request","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))

	// token auth
	tokenAuth := auth.MiddlewareWithConfig(auth.MiddlewareConfig{
		Key: key,
	})
	// basic auth
	basicAuth := auth.BasicAuthMiddleware(auth.BasicAuthMiddlewareConfig{
		Key: key,
	})

	csrfProtection := csrf.Middleware()

	// ----------------------------------------------------------------
	// handlers
	// ----------------------------------------------------------------
	e.GET("/", handlers.RootHandler)

	// actions endpoint
	e.POST("/actions/:name", handlers.FetchActionHandler(r, aFactory), tokenAuth)

	// "/api/..." endpoints are used to communicate with the client.

	// notify action result
	e.POST("/api/notify", handlers.NotifyActionResultHandler(r), tokenAuth)
	// session
	e.POST("/api/session/new", handlers.SessionNewHandler(cfg, r), tokenAuth)
	e.GET("/api/session/connect/:client_id/:session_id", handlers.SessionConnectHandler(r), tokenAuth)

	// token
	if cfg.ExposeNewToken {
		// api endpoint
		e.POST("/api/new-token", handlers.NewTokenHandler(tokenGenerator))
		// web ui
		e.GET("/new-token", handlers.NewTokenPageHandler(cfg), csrfProtection)
		e.POST("/new-token/create", handlers.NewTokenCreateHandler(cfg, tokenGenerator), csrfProtection)
	}

	// health check endpoint
	// It is useful if you deploy the server by using Kamal.
	// see: https://kamal-deploy.org/docs/configuration/proxy/#healthcheck
	e.GET("/up", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	// OpenAPI docs
	e.GET("/docs*", handlers.DocsHandler(r), basicAuth)

	// static files
	e.StaticFS("/", echo.MustSubFS(publicFS, "public"))

	// https://echo.labstack.com/docs/cookbook/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// start server
	go func() {
		if err := e.Start(e.Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Errorf("the server returned an error: %+v", err)
		}
	}()

	e.Logger.Infof("server started on %s", e.Server.Addr)

	// Wait for interrupt signal to stop the process.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("failed to shutdown the server: %+v", err)
	}

	return nil
}

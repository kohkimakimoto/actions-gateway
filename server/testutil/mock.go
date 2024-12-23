package testutil

//go:generate ../../../.dev/go-tools/bin/mockgen -destination=mock/logger.go -package=mock_logger github.com/labstack/echo/v4 Logger

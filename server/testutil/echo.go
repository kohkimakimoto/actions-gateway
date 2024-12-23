package testutil

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"testing"
)

func NewEchoInstance(t *testing.T) *echo.Echo {
	t.Helper()

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Debug = true
	e.Logger = NoopLogger
	return e
}

// NoopLogger is a logger that does nothing.
var NoopLogger = &NoopLoggerImpl{}

type NoopLoggerImpl struct {
}

func (l *NoopLoggerImpl) Output() io.Writer {
	return nil
}

func (l *NoopLoggerImpl) SetOutput(w io.Writer) {
}

func (l *NoopLoggerImpl) Prefix() string {
	return ""
}

func (l *NoopLoggerImpl) SetPrefix(p string) {
}

func (l *NoopLoggerImpl) Level() log.Lvl {
	return log.OFF
}

func (l *NoopLoggerImpl) SetLevel(v log.Lvl) {
}

func (l *NoopLoggerImpl) SetHeader(h string) {
}

func (l *NoopLoggerImpl) Print(i ...interface{}) {
}

func (l *NoopLoggerImpl) Printf(format string, args ...interface{}) {
}

func (l *NoopLoggerImpl) Printj(j log.JSON) {
}

func (l *NoopLoggerImpl) Debug(i ...interface{}) {
}

func (l *NoopLoggerImpl) Debugf(format string, args ...interface{}) {
}

func (l *NoopLoggerImpl) Debugj(j log.JSON) {
}

func (l *NoopLoggerImpl) Info(i ...interface{}) {
}

func (l *NoopLoggerImpl) Infof(format string, args ...interface{}) {
}

func (l *NoopLoggerImpl) Infoj(j log.JSON) {
}

func (l *NoopLoggerImpl) Warn(i ...interface{}) {
}

func (l *NoopLoggerImpl) Warnf(format string, args ...interface{}) {
}

func (l *NoopLoggerImpl) Warnj(j log.JSON) {
}

func (l *NoopLoggerImpl) Error(i ...interface{}) {
}

func (l *NoopLoggerImpl) Errorf(format string, args ...interface{}) {
}

func (l *NoopLoggerImpl) Errorj(j log.JSON) {
}

func (l *NoopLoggerImpl) Fatal(i ...interface{}) {
}

func (l *NoopLoggerImpl) Fatalj(j log.JSON) {
}

func (l *NoopLoggerImpl) Fatalf(format string, args ...interface{}) {
}

func (l *NoopLoggerImpl) Panic(i ...interface{}) {
}

func (l *NoopLoggerImpl) Panicf(format string, args ...interface{}) {
}

func (l *NoopLoggerImpl) Panicj(j log.JSON) {
}

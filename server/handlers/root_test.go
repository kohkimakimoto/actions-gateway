package handlers

import (
	"github.com/kohkimakimoto/actions-gateway/server/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	e := testutil.NewEchoInstance(t)
	e.GET("/", RootHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Actions Gateway is running. version: 0.0.0 (hash: unknown)", rec.Body.String())
}

package handlers

import (
	"encoding/json"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/kohkimakimoto/actions-gateway/server/testutil"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTokenHandler(t *testing.T) {
	e := testutil.NewEchoInstance(t)
	k, _ := auth.LoadKeyString("12345678901234567890123456789012")
	g := auth.NewTokenGenerator(k, auth.WithTime(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)), auth.WithClientId("00000000-0000-0000-0000-000000000001"))
	e.POST("/new-token", NewTokenHandler(g))

	req := httptest.NewRequest(http.MethodPost, "/new-token", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := &types.NewTokenResponse{}
	err := json.Unmarshal(rec.Body.Bytes(), resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2MDk0NTkyMDAsInN1YiI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMSJ9.33F9jqdaYWuJWiz3w68f7SwhVHZ3fqkqgQ1ofQtQ1bY", resp.Token)
}

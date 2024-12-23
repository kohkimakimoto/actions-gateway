package client

import (
	"bytes"
	"github.com/kohkimakimoto/actions-gateway/client/config"
	"github.com/kohkimakimoto/actions-gateway/server/types"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestClient_makeURL(t *testing.T) {
	testCases := []struct {
		name     string
		client   *Client
		urlPath  string
		expected string
	}{
		{
			name: "server URL has a trailing slash",
			client: &Client{
				config: &config.Config{
					Server: "http://localhost:8080/",
				},
			},
			urlPath:  "/api/new-token",
			expected: "http://localhost:8080/api/new-token",
		},
		{
			name: "server URL does not have a trailing slash",
			client: &Client{
				config: &config.Config{
					Server: "http://localhost:8080",
				},
			},
			urlPath:  "/api/new-token",
			expected: "http://localhost:8080/api/new-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.client.makeURL(tc.urlPath)
			if actual != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, actual)
			}
		})
	}
}

func TestClient_NewToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := New(&config.Config{
			Server: "http://localhost:8080",
		}, nil, nil)
		client.httpClient = testHttpClient(t, func(req *http.Request) *http.Response {
			// Check the request
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "http://localhost:8080/api/new-token", req.URL.String())

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"token": "this-is-a-test-token"}`)),
			}
		})

		token, err := client.NewToken()
		assert.NoError(t, err)
		assert.Equal(t, "this-is-a-test-token", token)
	})
}

func TestClient_NotifyResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := New(&config.Config{
			Server: "http://localhost:8080",
		}, nil, nil)
		client.httpClient = testHttpClient(t, func(req *http.Request) *http.Response {
			// Check the request
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "http://localhost:8080/api/notify", req.URL.String())

			return &http.Response{
				StatusCode: http.StatusNoContent,
			}
		})

		err := client.NotifyResult(&types.ActionResult{
			Status: types.ActionResultStatusSuccess,
		})
		assert.NoError(t, err)
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testHttpClient(t *testing.T, fn RoundTripFunc) *http.Client {
	t.Helper()
	return &http.Client{
		Transport: fn,
	}
}

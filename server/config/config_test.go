package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_WebSocketURL(t *testing.T) {
	testCases := map[string]struct {
		cfg      *Config
		expected string
	}{
		"http": {
			cfg: &Config{
				URL: "http://localhost:8080",
			},
			expected: "ws://localhost:8080",
		},
		"https": {
			cfg: &Config{
				URL: "https://localhost:8080",
			},
			expected: "wss://localhost:8080",
		},
		"no scheme": {
			cfg: &Config{
				URL: "localhost:8080",
			},
			expected: "ws://localhost:8080",
		},
		"trailing slash": {
			cfg: &Config{
				URL: "http://localhost:8080/",
			},
			expected: "ws://localhost:8080",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.cfg.WebSocketURL())
		})
	}
}

func TestUpdateByFile(t *testing.T) {
	t.Run("use config from file", func(t *testing.T) {
		cfg := New()
		f := testTempFile(t, []byte(`
addr = ":8080"
url = "http://localhost:8080"
secret = "test_secret"
expose_new_token = true
debug = true
`))
		err := UpdateByFile(cfg, f.Name())
		assert.NoError(t, err)
		assert.Equal(t, ":8080", cfg.Addr)
		assert.Equal(t, "http://localhost:8080", cfg.URL)
		assert.Equal(t, "test_secret", cfg.Secret)
		assert.True(t, cfg.ExposeNewToken)
		assert.True(t, cfg.Debug)
	})

	t.Run("use default config", func(t *testing.T) {
		cfg := New()
		f := testTempFile(t, []byte(``))
		err := UpdateByFile(cfg, f.Name())
		assert.NoError(t, err)
		assert.Equal(t, ":18800", cfg.Addr)
		assert.Equal(t, "http://localhost:18800", cfg.URL)
		assert.Equal(t, "", cfg.Secret)
		assert.False(t, cfg.ExposeNewToken)
		assert.False(t, cfg.Debug)
	})

	t.Run("fail to load config from file", func(t *testing.T) {
		cfg := New()
		f := testTempFile(t, []byte(`invalid`))
		err := UpdateByFile(cfg, f.Name())
		assert.Error(t, err)
	})
}

func TestUpdateByEnvironments(t *testing.T) {
	t.Run("update config from environment variables", func(t *testing.T) {
		_ = os.Setenv("ACTIONS_GATEWAY_ADDR", ":8080")
		_ = os.Setenv("ACTIONS_GATEWAY_URL", "http://localhost:8080")
		_ = os.Setenv("ACTIONS_GATEWAY_SECRET", "test_secret")
		_ = os.Setenv("ACTIONS_GATEWAY_EXPOSE_NEW_TOKEN", "true")
		_ = os.Setenv("ACTIONS_GATEWAY_DEBUG", "true")
		defer func() {
			_ = os.Unsetenv("ACTIONS_GATEWAY_ADDR")
			_ = os.Unsetenv("ACTIONS_GATEWAY_URL")
			_ = os.Unsetenv("ACTIONS_GATEWAY_SECRET")
			_ = os.Unsetenv("ACTIONS_GATEWAY_EXPOSE_NEW_TOKEN")
			_ = os.Unsetenv("ACTIONS_GATEWAY_DEBUG")
		}()
		cfg := New()
		UpdateByEnvironments(cfg)
		assert.Equal(t, ":8080", cfg.Addr)
		assert.Equal(t, "http://localhost:8080", cfg.URL)
		assert.Equal(t, "test_secret", cfg.Secret)
		assert.True(t, cfg.ExposeNewToken)
		assert.True(t, cfg.Debug)
	})
}

func testTempFile(t *testing.T, b []byte) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.Name(), b, 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	})
	return f
}

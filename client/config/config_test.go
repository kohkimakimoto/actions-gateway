package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_Dir(t *testing.T) {
	c := &Config{
		Path: "/path/to/config.toml",
	}
	assert.Equal(t, "/path/to", c.Dir())
}

func TestLoadFromFile(t *testing.T) {
	t.Run("use config from file", func(t *testing.T) {
		f := testTempFile(t, []byte(`
server = "http://localhost:8080"
`))
		cfg, err := LoadFromFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost:8080", cfg.Server)
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

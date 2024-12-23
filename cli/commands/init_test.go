package commands

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCommand(t *testing.T) {
	app := cli.NewApp()
	app.Writer = &bytes.Buffer{}
	app.Commands = []*cli.Command{
		InitCommand,
	}

	dir := filepath.Join(testTempDir(t), "config")
	// t.Logf(dir)

	err := app.Run([]string{"", "init", "-d", dir})
	assert.NoError(t, err)

	b, err := os.ReadFile(filepath.Join(dir, "config.toml"))
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
}

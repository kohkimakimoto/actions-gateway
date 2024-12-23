package commands

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"os"
	"testing"
)

func TestGojqCommand(t *testing.T) {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		GojqCommand,
	}

	// gojq uses STDIN and STDOUT directly in the command.
	// So, we need to mock them.

	// mock stdin
	oStdin := os.Stdin
	defer func() {
		os.Stdin = oStdin
	}()
	rIn, wIn, _ := os.Pipe()
	os.Stdin = rIn

	// mock stdout
	oStdout := os.Stdout
	defer func() {
		os.Stdout = oStdout
	}()
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// input test data to stdin
	_, _ = wIn.WriteString(`{"key": "value"}`)
	_ = wIn.Close()

	err := app.Run([]string{"", "gojq", ".key"})
	assert.NoError(t, err)

	_ = wOut.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(rOut)
	assert.NoError(t, err)

	assert.Equal(t, "\"value\"\n", buf.String())
}

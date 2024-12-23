package commands

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpecCommand(t *testing.T) {
	t.Run("output spec with initial config", func(t *testing.T) {
		// init config
		app := cli.NewApp()
		app.Writer = &bytes.Buffer{}
		app.Commands = []*cli.Command{
			InitCommand,
		}
		dir := filepath.Join(testTempDir(t), "config")
		err := app.Run([]string{"", "init", "-d", dir})
		assert.NoError(t, err)

		// run spec command
		app = cli.NewApp()
		outBuffer := &bytes.Buffer{}
		app.Writer = outBuffer
		app.Commands = []*cli.Command{
			SpecCommand,
		}

		err = app.Run([]string{"", "spec", "-c", filepath.Join(dir, "config.toml")})
		assert.NoError(t, err)
		spec := outBuffer.String()
		assert.Equal(t, strings.TrimLeft(`
openapi: 3.1.0
info:
  title: Actions Gateway API
  summary: Actions Gateway API
  description: Actions Gateway API
  version: 1.0.0
servers:
  - url: https://actions-gateway.kohkimakimoto.dev
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas: {}
security:
  - bearerAuth: []
paths:
  /actions/openURL:
    post:
      summary: Open a URL
      description: |
        This action receives a URL and opens it in the default browser.
      operationId: openURL
      x-openai-isConsequential: false
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                url:
                  type: string
                  description: The URL to open
                  example: "https://github.com"
              required:
                - url
      responses:
        "200":
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  opened_url:
                    type: string
                    description: The URL that was opened
                    example: "https://github.com"
                required:
                  - opened_url
`, "\n"), spec)
	})
}

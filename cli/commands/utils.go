package commands

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/kohkimakimoto/actions-gateway/client"
	"github.com/kohkimakimoto/actions-gateway/client/config"
	"github.com/urfave/cli/v2"
	"io"
	"regexp"
)

var clientConfigFlag = &cli.StringFlag{
	Name:    "config",
	Aliases: []string{"c"},
	Usage:   "Path to the client config `file`",
}

func getClientConfig(cCtx *cli.Context) (*config.Config, error) {
	configFIle := cCtx.String("config")
	if configFIle == "" {
		_configFile, err := config.DefaultFilepath()
		if err != nil {
			return nil, fmt.Errorf("failed to get default config: %w", err)
		}
		configFIle = _configFile
	}

	cfg, err := config.LoadFromFile(configFIle)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configFIle, err)
	}
	return cfg, nil
}

func newClient(cCtx *cli.Context, cfg *config.Config) *client.Client {
	return client.New(cfg, cCtx.App.Writer, cCtx.App.ErrWriter)
}

type SimpleTableWriter struct {
	table.Writer
	Out io.Writer
}

func newSimpleTableWriter(out io.Writer) *SimpleTableWriter {
	// custom borderless table style
	style := table.StyleDefault
	style.Box.PaddingLeft = ""
	style.Box.PaddingRight = ""
	style.Box.MiddleVertical = "   "
	style.Options = table.Options{
		DrawBorder:      false,
		SeparateColumns: true,
		SeparateFooter:  false,
		SeparateHeader:  false,
		SeparateRows:    false,
	}
	style.Format.Header = text.FormatDefault

	w := table.NewWriter()
	w.SetStyle(style)
	return &SimpleTableWriter{
		Writer: w,
		Out:    out,
	}
}

var reRemoveTrailingSpace = regexp.MustCompile(`\s+\n`)

func (t *SimpleTableWriter) Render() string {
	// Wrap the table.Writer's Render() method to remove trailing spaces.
	outStr := t.Writer.Render()
	outStr = reRemoveTrailingSpace.ReplaceAllString(outStr, "\n")
	_, _ = fmt.Fprintln(t.Out, outStr)
	return outStr
}

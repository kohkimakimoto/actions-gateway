package cli

import (
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/cli/commands"
	"github.com/kohkimakimoto/actions-gateway/version"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func Main(args []string) {
	app := cli.NewApp()
	app.Name = "actions-gateway"
	app.Version = version.Version
	app.Usage = "An API server that allows running local programs through HTTP requests."
	app.Commands = []*cli.Command{
		commands.GojqCommand,
		commands.InitCommand,
		commands.LogsCommand,
		commands.NewTokenCommand,
		commands.ServeCommand,
		commands.SpecCommand,
		commands.StartCommand,
		commands.StatusCommand,
		commands.StopCommand,
	}
	if err := app.Run(args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
}

func init() {
	cli.AppHelpTemplate = strings.TrimLeft(`
Usage: {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Args}} [arguments...]{{end}}{{end}}{{end}}

{{ .Usage }}{{if .VisibleCommands}}

Commands:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

Global options:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

Global options:{{template "visibleFlagTemplate" .}}{{end}}{{if .Description}}

Description:
   {{template "descriptionTemplate" .}}{{end}}
`, "\n")

	cli.CommandHelpTemplate = strings.TrimLeft(`
Usage: {{template "usageTemplate" .}}

{{ .Usage }}{{if .VisibleFlagCategories}}

Options:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

Options:{{template "visibleFlagTemplate" .}}{{end}}{{if .Description}}

Description:
   {{template "descriptionTemplate" .}}{{end}}
`, "\n")
}

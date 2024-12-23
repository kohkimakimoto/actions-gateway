package commands

import (
	"fmt"

	"github.com/kohkimakimoto/actions-gateway/server"
	"github.com/kohkimakimoto/actions-gateway/server/config"
	"github.com/urfave/cli/v2"
)

var ServeCommand = &cli.Command{
	Name:  "serve",
	Usage: "Start the server process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "server-config",
			Usage: "The configuration `file` for the server",
		},
	},
	Action: serveAction,
}

func serveAction(cCtx *cli.Context) error {
	cfg := config.New()

	configFIle := cCtx.String("server-config")
	if configFIle != "" {
		if err := config.UpdateByFile(cfg, configFIle); err != nil {
			return fmt.Errorf("failed to load config file %s: %w", configFIle, err)
		}
	}

	config.UpdateByEnvironments(cfg)

	if cfg.Secret == "" {
		return fmt.Errorf("secret is required. Please set ACTIONS_GATEWAY_SECRET or set 'secret' in the config file")
	}

	return server.Start(cfg)
}

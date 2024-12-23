package commands

import (
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/builtin"
	"github.com/kohkimakimoto/actions-gateway/client/config"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "Initialize a client config directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "directory",
			Aliases: []string{"d"},
			Usage:   "Specify a `directory` to init as a client config directory",
		},
	},
	Action: initAction,
}

func initAction(cCtx *cli.Context) error {
	dir := cCtx.String("directory")
	if dir == "" {
		_dir, err := config.DefaultDir()
		if err != nil {
			return err
		}
		dir = _dir
	}

	// create config directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.FileMode(0700)); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	} else {
		return fmt.Errorf("directory %s already exists", dir)
	}

	// create actions directory
	actionsDir := filepath.Join(dir, "actions")
	if err := os.MkdirAll(actionsDir, os.FileMode(0700)); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", actionsDir, err)
	}

	// create builtin actions
	if err := builtin.InitFiles(actionsDir); err != nil {
		return fmt.Errorf("failed to create builtin actions: %w", err)
	}

	// create config file
	configFile := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(configFile, []byte(config.InitialConfig), os.FileMode(0600)); err != nil {
		return fmt.Errorf("failed to create config file %s: %w", configFile, err)
	}

	_, _ = fmt.Fprintln(cCtx.App.Writer, "created: "+dir)
	return nil
}

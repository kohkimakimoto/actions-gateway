package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"os/exec"
)

var LogsCommand = &cli.Command{
	Name:  "logs",
	Usage: "Show the logs of the client agent",
	Flags: []cli.Flag{
		clientConfigFlag,
		&cli.BoolFlag{
			Name:    "follow",
			Aliases: []string{"f"},
			Usage:   "Follow the log file",
		},
	},
	Action: logsAction,
}

func logsAction(cCtx *cli.Context) error {
	config, err := getClientConfig(cCtx)
	if err != nil {
		return err
	}

	if config.LogFile == "" {
		return fmt.Errorf("specifying log_file is required in your config file to show logs")
	}

	if cCtx.Bool("follow") {
		// TODO: reimplement "tail" operation in Go code without using exec.Command
		cmd := exec.Command("tail", "-f", config.LogAbsFile)
		cmd.Stdout = cCtx.App.Writer
		cmd.Stderr = cCtx.App.ErrWriter
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run tail command: %w", err)
		}
	} else {
		file, err := os.Open(config.LogAbsFile)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("no log file was found: %w", err)
			}
			return fmt.Errorf("failed to read log file %s: %w", config.LogAbsFile, err)
		}
		defer file.Close()

		if _, err = io.Copy(cCtx.App.Writer, file); err != nil {
			return fmt.Errorf("failed to show logs: %w", err)
		}
	}

	return nil
}

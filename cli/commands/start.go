package commands

import (
	"context"
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/client/actions"
	"github.com/kohkimakimoto/actions-gateway/client/status"
	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v2"
	"os/signal"
	"syscall"
)

var StartCommand = &cli.Command{
	Name:  "start",
	Usage: "Start the client agent to connect to the server",
	Flags: []cli.Flag{
		clientConfigFlag,
		&cli.BoolFlag{
			Name:    "daemon",
			Aliases: []string{"d"},
			Usage:   "Run as a daemon (background) process",
		},
	},
	Action: startAction,
}

func startAction(cCtx *cli.Context) (err error) {
	config, err := getClientConfig(cCtx)
	if err != nil {
		return err
	}

	if cCtx.Bool("daemon") {
		if config.StatusFile == "" {
			return fmt.Errorf("specifying status_file is required in your config file to run as a daemon")
		}
		if config.PidFile == "" {
			return fmt.Errorf("specifying pid_file is required in your config file to run as a daemon")
		}
		if config.LogFile == "" {
			return fmt.Errorf("specifying log_file is required in your config file to run as a daemon")
		}

		d := daemon.Context{
			PidFileName: config.PidAbsFile,
			PidFilePerm: 0644,
			LogFileName: config.LogAbsFile,
			LogFilePerm: 0644,
		}

		child, err := d.Reborn()
		if err != nil {
			return fmt.Errorf("failed to run as a daemon: %w", err)
		}
		if child != nil {
			// If a child process object is returned from Reborn(), it indicates that the current process is the parent process.
			// So, it should exit.
			return nil
		}
		defer func() {
			_ = d.Release()
		}()
	}

	// init the action manager
	am, err := actions.NewActionManager(config)
	if err != nil {
		return err
	}

	// init the status writer
	sw := status.NewWriter(config.StatusAbsFile)
	if err := sw.Init(); err != nil {
		return err
	}
	defer func() {
		_ = sw.UpdateToInactive(err)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// connect to the server
	go func() {
		if err2 := newClient(cCtx, config).Connect(am, sw); err2 != nil {
			err = fmt.Errorf("failed to connect to the server: %w", err2)
		}
		stop()
	}()

	// wait for interrupt signal
	<-ctx.Done()
	return err
}

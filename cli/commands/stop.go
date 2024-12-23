package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"syscall"
)

var StopCommand = &cli.Command{
	Name:  "stop",
	Usage: "Stop the client agent daemon process",
	Flags: []cli.Flag{
		clientConfigFlag,
	},
	Action: disconnectAction,
}

func disconnectAction(cCtx *cli.Context) error {
	config, err := getClientConfig(cCtx)
	if err != nil {
		return err
	}

	if config.PidFile == "" {
		return fmt.Errorf("specifying pid_file is required in your config file to disconnect")
	}

	pidData, err := os.ReadFile(config.PidAbsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no daemon process was found (pid file %s not found): %w", config.PidAbsFile, err)
		}
		return fmt.Errorf("failed to read pid file %s: %w", config.PidAbsFile, err)
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		return fmt.Errorf("failed to parse pid file %s: %w", config.PidAbsFile, err)
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}
	err = p.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to send SIGTERM to process %d: %w", pid, err)
	}
	return nil
}

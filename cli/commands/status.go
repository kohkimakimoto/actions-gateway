package commands

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kohkimakimoto/actions-gateway/client/status"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
)

var StatusCommand = &cli.Command{
	Name:  "status",
	Usage: "Show the status of the client agent",
	Flags: []cli.Flag{
		clientConfigFlag,
	},
	Action: statusAction,
}

func statusAction(cCtx *cli.Context) error {
	config, err := getClientConfig(cCtx)
	if err != nil {
		return err
	}

	// Read the PID file
	var pid string
	if config.PidFile != "" {
		pidData, err := os.ReadFile(config.PidAbsFile)
		if err != nil {
			if os.IsNotExist(err) {
				pid = "no daemon process was found"
			} else {
				pid = fmt.Sprintf("failed to read pid file %s", config.PidAbsFile)
			}
		} else {
			_pid, err := strconv.Atoi(string(pidData))
			if err != nil {
				pid = fmt.Sprintf("failed to parse pid file %s", config.PidAbsFile)
			} else {
				pid = strconv.Itoa(_pid)
			}
		}
	}

	// Read the status file
	r := status.NewReader(config.StatusAbsFile)
	s, err := r.Read()
	if err != nil {
		return err
	}

	t := newSimpleTableWriter(cCtx.App.Writer)
	t.AppendRow(table.Row{"Pid:", pid})
	t.AppendRow(table.Row{"Status:", s.StatusCode})
	t.Render()

	return nil
}

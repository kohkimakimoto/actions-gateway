package commands

import (
	gojqcli "github.com/itchyny/gojq/cli"
	"github.com/urfave/cli/v2"
	"os"
)

var GojqCommand = &cli.Command{
	Name:            "gojq",
	Usage:           "built-in gojq command",
	SkipFlagParsing: true,
	Action:          gojqAction,
}

func gojqAction(cCtx *cli.Context) error {
	// override os.Args to run gojq inside actions-gateway
	os.Args = append([]string{"actions-gateway-gojq"}, cCtx.Args().Slice()...)
	if exitStatus := gojqcli.Run(); exitStatus != 0 {
		return cli.Exit("", exitStatus)
	}
	return nil
}

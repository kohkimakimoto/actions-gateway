package commands

import (
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/client/actions"
	"github.com/urfave/cli/v2"
)

var SpecCommand = &cli.Command{
	Name:  "spec",
	Usage: "Output the OpenAPI spec of your actions",
	Flags: []cli.Flag{
		clientConfigFlag,
	},
	Action: specAction,
}

func specAction(cCtx *cli.Context) error {
	config, err := getClientConfig(cCtx)
	if err != nil {
		return err
	}

	am, err := actions.NewActionManager(config)
	if err != nil {
		return err
	}

	spec, err := am.OutputSpec(cCtx.App.ErrWriter)
	if err != nil {
		return fmt.Errorf("failed to output OpenAPI spec: %w", err)
	}

	_, _ = fmt.Fprint(cCtx.App.Writer, spec)
	return nil
}

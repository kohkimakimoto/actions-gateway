package commands

import (
	"fmt"
	"github.com/kohkimakimoto/actions-gateway/server/auth"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"syscall"
)

var NewTokenCommand = &cli.Command{
	Name:   "new-token",
	Usage:  "Generate a new token",
	Action: newTokenAction,
	Flags: []cli.Flag{
		clientConfigFlag,
		&cli.BoolFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "Generate a new token without connecting to the server",
		},
	},
}

func newTokenAction(cCtx *cli.Context) error {
	if cCtx.Bool("local") {
		return newToken(cCtx)
	} else {
		return newTokenUsingServer(cCtx)
	}
}

func newToken(cCtx *cli.Context) error {
	_, _ = fmt.Fprint(cCtx.App.Writer, "Enter your secret: ")
	bSecret, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read the secret: %w", err)
	}
	_, _ = fmt.Fprintln(cCtx.App.Writer)
	key, err := auth.LoadKey(bSecret)
	if err != nil {
		return fmt.Errorf("failed to load the secret: %w", err)
	}

	token, err := auth.NewTokenGenerator(key).NewTokenAsJWTString()
	if err != nil {
		return fmt.Errorf("failed to generate a new token: %w", err)
	}
	_, _ = fmt.Fprintln(cCtx.App.Writer, token)
	return nil
}

func newTokenUsingServer(cCtx *cli.Context) error {
	config, err := getClientConfig(cCtx)
	if err != nil {
		return err
	}

	token, err := newClient(cCtx, config).NewToken()
	if err != nil {
		return fmt.Errorf("failed to generate a new token: %w", err)
	}
	_, _ = fmt.Fprintln(cCtx.App.Writer, token)
	return nil
}

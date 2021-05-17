package app

import (
	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/appflag"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if c.Args().Present() {
		cli.ShowSubcommandHelp(c)
		return exit.Status(1)
	}

	if appflag.Version() {
		cli.ShowVersion(c)
		return nil
	}

	return cli.ShowAppHelp(c)
}

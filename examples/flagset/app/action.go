package app

import (
	"fmt"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/examples/flagset/app/flagapp"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if c.Args().Present() {
		_ = cli.ShowAppHelp(c)
		return exit.Status(1)
	}

	if v, ok := flagapp.LookupCount(); ok {
		fmt.Fprintln(c.App.Writer, "count", v)
	} else {
		fmt.Fprintln(c.App.Writer, "count not set")
	}
	return nil
}

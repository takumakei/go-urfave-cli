package split

import (
	"errors"
	"fmt"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/examples/flagset/app/flagapp"
	"github.com/takumakei/go-urfave-cli/examples/flagset/app/split/flagsplit"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if c.Args().Present() {
		_ = cli.ShowSubcommandHelp(c)
		return exit.Status(1)
	}

	if v, ok := flagapp.LookupCount(); ok {
		fmt.Fprintln(c.App.Writer, "count", v)
	} else {
		fmt.Fprintln(c.App.Writer, "count not set")
	}

	err := flagsplit.Direction.Select(map[cli.Flag]func() error{
		flagsplit.FlagHorizontal: func() error {
			fmt.Fprintln(c.App.Writer, "horizontal", flagsplit.Horizontal())
			return nil
		},

		flagsplit.FlagVertical: func() error {
			fmt.Fprintln(c.App.Writer, "vertical", flagsplit.Vertical())
			return nil
		},

		nil: func() error {
			return exit.Error(2, errors.New("-horizontal or -vertical is required"))
		},
	})

	return err
}

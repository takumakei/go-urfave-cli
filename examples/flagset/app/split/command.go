package split

import (
	"github.com/takumakei/go-urfave-cli/examples/flagset/app/split/flagsplit"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:      "split",
	Usage:     "split window",
	ArgsUsage: " ",
	Flags:     flagsplit.Flags,
	Before:    flagsplit.FlagSet.Init,
	Action:    Action,
}

package client

import (
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/client/clientflag"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:      "client",
	Usage:     "client",
	ArgsUsage: " ",
	Flags:     clientflag.Flags,
	Before:    clientflag.Before,
	Action:    Action,
}

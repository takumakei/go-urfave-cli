package server

import (
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/server/serverflag"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:      "server",
	Usage:     "server",
	ArgsUsage: " ",
	Flags:     serverflag.Flags,
	Before:    serverflag.Before,
	Action:    Action,
}

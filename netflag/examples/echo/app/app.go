package app

import (
	"github.com/takumakei/go-urfave-cli/fish"
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/appflag"
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/client"
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/server"
	"github.com/urfave/cli/v2"
)

var App = cli.NewApp()

func init() {
	App.Usage = "echo"
	App.ArgsUsage = " "
	App.HideVersion = true
	App.Flags = appflag.Flags
	App.Before = appflag.FlagSet.Init
	App.Action = Action
	App.Commands = []*cli.Command{
		server.Command,
		client.Command,
	}
	App.EnableBashCompletion = true
	fish.EnableFishCompletionCommand(App)
}

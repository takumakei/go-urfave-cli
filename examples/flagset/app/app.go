package app

import (
	"github.com/takumakei/go-urfave-cli/examples/flagset/app/flagapp"
	"github.com/takumakei/go-urfave-cli/examples/flagset/app/split"
	"github.com/takumakei/go-urfave-cli/fish"
	"github.com/urfave/cli/v2"
)

var App = cli.NewApp()

func init() {
	App.Usage = "example cli application"
	App.ArgsUsage = " "
	App.Flags = flagapp.Flags
	App.Before = flagapp.Before
	App.After = flagapp.After
	App.Action = Action
	App.Commands = []*cli.Command{
		split.Command,
	}
	App.EnableBashCompletion = true

	fish.EnableFishCompletionCommand(App)
}

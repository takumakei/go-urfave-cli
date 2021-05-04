package main

import (
	"os"

	"github.com/takumakei/go-urfave-cli/fish"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	fish.EnableFishCompletionCommand(app)
	_ = app.Run(os.Args)
}

package main

import (
	"os"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/zapflag"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	zf := zapflag.New(clix.FlagPrefix("MINI_"))

	app := cli.NewApp()
	app.Flags = zf.Flags()
	app.Before = zf.InitGlobal
	app.After = func(c *cli.Context) error {
		zap.L().Sync()
		return nil
	}
	app.Action = func(c *cli.Context) error {
		zap.L().Info("hello world")
		return nil
	}
	exit.Exit(app.Run(os.Args))
}

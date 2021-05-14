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
	fp := clix.FlagPrefix("NEWNAME_")

	zfGlobal := zapflag.New(fp)
	zfSecond := zapflag.NewName("second", fp)

	var second *zap.Logger

	app := cli.NewApp()
	app.Flags = clix.Flags(zfGlobal.Flags(), zfSecond.Flags())
	app.Before = clix.Chain(zfGlobal.InitGlobal, zfSecond.InitInto(&second))
	app.After = clix.Chain(
		zapflag.Sync(&second, zapflag.IgnoreError),
		zapflag.SyncGlobal(zapflag.IgnoreError),
	)
	app.Action = func(c *cli.Context) error {
		zap.L().Info("hello world")
		second.Info("hello world")
		return nil
	}
	exit.Exit(app.Run(os.Args))
}

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

	zf := zapflag.New(fp)

	zf2 := zapflag.NewName("zap", fp)
	var log2 *zap.Logger
	initLog2 := func(c *cli.Context) error {
		v, err := zf2.Logger()
		if err == nil {
			log2 = v.With(zap.String("logger", "zap"))
		}
		return err
	}

	app := cli.NewApp()
	app.Flags = clix.Flags(zf.Flags(), zf2.Flags())
	app.Before = clix.Chain(zf.InitGlobal, zf2.Init, initLog2)
	app.After = func(c *cli.Context) error {
		zap.L().Sync()
		log2.Sync()
		return nil
	}
	app.Action = func(c *cli.Context) error {
		zap.L().Info("hello world")
		log2.Info("hello world")
		return nil
	}
	exit.Exit(app.Run(os.Args))
}

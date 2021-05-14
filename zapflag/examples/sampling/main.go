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
	fp := clix.FlagPrefix("SAMPLING_")

	zf := zapflag.New(fp)

	flagCount := &cli.IntFlag{
		Name:        "count",
		Aliases:     []string{"n"},
		Usage:       "repeat `count`",
		EnvVars:     fp.EnvVars("COUNT", "N"),
		FilePath:    fp.FilePath("COUNT", "N"),
		Value:       7,
		Destination: new(int),
	}

	flagWorld := &cli.BoolFlag{
		Name:        "world",
		Aliases:     []string{"w"},
		Usage:       "log world",
		EnvVars:     fp.EnvVars("WORLD", "W"),
		FilePath:    fp.FilePath("WORLD", "W"),
		Value:       true,
		Destination: new(bool),
	}

	app := cli.NewApp()
	app.Flags = clix.Flags(flagCount, flagWorld, zf.Flags())
	app.Before = zf.InitGlobal
	app.After = zapflag.SyncGlobal(zapflag.IgnoreError)
	app.Action = func(c *cli.Context) error {
		count := *flagCount.Destination
		world := *flagWorld.Destination

		for i := 0; i < count; i++ {
			zap.L().Info("hello", zap.Int("i", i))
			if world {
				zap.L().Info("world", zap.Int("i", i))
			}
		}
		return nil
	}

	exit.Exit(app.Run(os.Args))
}

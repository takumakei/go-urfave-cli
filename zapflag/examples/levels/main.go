package main

import (
	"os"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/zapflag"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	fp := clix.FlagPrefix("LEVELS_")

	zf := zapflag.New(fp)

	flagDo := &cli.GenericFlag{
		Name:     "do",
		Usage:    "do `level` [dpanic|panic|fatal]",
		EnvVars:  fp.EnvVars("DO"),
		FilePath: fp.FilePath("DO"),
		Value:    new(zapcore.Level),
	}

	flagSet := clix.NewFlagSet()

	app := cli.NewApp()
	app.Flags = clix.Flags(flagDo, zf.Flags())
	app.Before = clix.Chain(flagSet.Init, zf.InitGlobal)
	app.After = zapflag.SyncGlobal
	app.Action = func(c *cli.Context) error {
		zap.L().Debug("debug message")
		zap.L().Info("info message")
		zap.L().Warn("warn message")
		zap.L().Error("error message")

		if flagSet.IsSet(flagDo) {
			switch v := *flagDo.Value.(*zapcore.Level); v {
			case zap.DPanicLevel:
				zap.L().DPanic("dpanic message")
			case zap.PanicLevel:
				zap.L().Panic("panic message")
			case zap.FatalLevel:
				zap.L().Fatal("fatal message")
			default:
				zap.L().Warn("ignore --do", zap.String("level", v.String()))
			}
		}

		return nil
	}

	exit.Exit(app.Run(os.Args))
}

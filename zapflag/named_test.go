package zapflag_test

import (
	"github.com/takumakei/go-urfave-cli/zapflag"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func ExampleNamed() {
	log := zap.NewExample()

	app := cli.NewApp()
	app.Name = "example"
	app.Action = func(c *cli.Context) error { return nil }
	app.Before = func(c *cli.Context) error {
		log := zapflag.Named(log, c)
		log.Info("named logger for root")
		return nil
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:   "hello",
			Action: func(c *cli.Context) error { return nil },
			Before: func(c *cli.Context) error {
				log := zapflag.Named(log, c)
				log.Info("named logger for hello")
				return nil
			},
			Subcommands: []*cli.Command{
				&cli.Command{
					Name:   "world",
					Action: func(c *cli.Context) error { return nil },
					Before: func(c *cli.Context) error {
						log := zapflag.Named(log, c)
						log.Info("named logger for world")
						return nil
					},
					Subcommands: []*cli.Command{
						&cli.Command{
							Name:   "melon",
							Action: func(c *cli.Context) error { return nil },
							Before: func(c *cli.Context) error {
								log := zapflag.Named(log, c)
								log.Info("named logger for melon")
								return nil
							},
						},
					},
				},
			},
		},
	}
	_ = app.Run([]string{"example", "hello", "world", "melon"})
	// output:
	// {"level":"info","logger":"example","msg":"named logger for root"}
	// {"level":"info","logger":"example.hello","msg":"named logger for hello"}
	// {"level":"info","logger":"example.hello.world","msg":"named logger for world"}
	// {"level":"info","logger":"example.hello.world.melon","msg":"named logger for melon"}
}

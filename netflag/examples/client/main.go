package main

import (
	"errors"
	"io"
	"os"

	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/netflag"
	"github.com/urfave/cli/v2"
)

func main() {
	client := netflag.NewClient(clix.FlagPrefix("CLIENT_"))

	app := cli.NewApp()
	app.Flags = client.Flags()
	app.Before = client.Before
	app.Action = func(c *cli.Context) error {
		conn, err := client.Dial()
		if err != nil {
			return err
		}
		defer delint.AnywayFunc(conn.Close)

		p := make([]byte, 2048)
		for {
			n, err := conn.Read(p)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
			c.App.Writer.Write(p[:n])
		}
	}
	exit.Exit(app.Run(os.Args))
}

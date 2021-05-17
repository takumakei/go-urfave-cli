package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/netflag"
	"github.com/urfave/cli/v2"
)

func main() {
	server := netflag.NewServer(clix.FlagPrefix("SERVER_"))

	app := cli.NewApp()
	app.Flags = server.Flags()
	app.Before = server.Before
	app.Action = func(c *cli.Context) error {
		lis, err := server.Listen()
		if err != nil {
			return err
		}
		defer delint.AnywayFunc(lis.Close)

		for {
			conn, err := lis.Accept()
			if err != nil {
				fmt.Fprintln(c.App.ErrWriter, err)
				conn.Close()
				continue
			}
			_, err = conn.Write([]byte(time.Now().Format(time.RFC3339Nano) + "\n"))
			if err != nil && !errors.Is(err, io.EOF) {
				fmt.Fprintln(c.App.ErrWriter, err)
			}
			err = conn.Close()
			if err != nil {
				fmt.Fprintln(c.App.ErrWriter, err)
			}
		}
	}
	exit.Exit(app.Run(os.Args))
}

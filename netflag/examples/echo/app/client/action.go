package client

import (
	"fmt"
	"time"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/client/clientflag"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if c.Args().Present() {
		cli.ShowSubcommandHelp(c)
		return exit.Status(1)
	}

	conn, err := clientflag.FlagConnect.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		for {
			p := make([]byte, 256)
			n, err := conn.Read(p)
			if err != nil {
				return
			}
			fmt.Printf("recv[%s]\n", string(p[:n]))
		}
	}()

	j := clientflag.Count()
	for i := 0; i < j; i++ {
		if _, err := conn.Write([]byte("hello world")); err != nil {
			return err
		}
	}

	time.Sleep(time.Second)

	return nil
}

package server

import (
	"fmt"
	"net"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app/server/serverflag"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if c.Args().Present() {
		cli.ShowSubcommandHelp(c)
		return exit.Status(1)
	}

	lis, err := serverflag.FlagListen.Listen()
	if err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	p := make([]byte, 256)
	for {
		n, err := conn.Read(p)
		if err != nil {
			return
		}

		fmt.Printf("recv[%s]\n", string(p[:n]))

		j := serverflag.Count()
		for i := 0; i < j; i++ {
			_, err = conn.Write(p[:n])
			if err != nil {
				return
			}
		}
	}
}

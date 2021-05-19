package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/grpcflag"
	"github.com/urfave/cli/v2"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var flag = grpcflag.NewDialer(
	clix.FlagPrefix("GREETER_"),
	grpcflag.Address("127.0.0.1:9900"),
)

func main() {
	app := cli.NewApp()
	app.Usage = "greeter client"
	app.Flags = flag.Flags()
	app.Before = flag.Before
	app.Action = Action
	exit.Exit(app.Run(os.Args))
}

func Action(c *cli.Context) error {
	ctx, cancel := context.WithTimeout(c.Context, 3*time.Second)
	conn, err := flag.DialContext(ctx)
	cancel()
	if err != nil {
		return err
	}

	greeter := pb.NewGreeterClient(conn)

	ctx, cancel = context.WithTimeout(c.Context, 3*time.Second)
	reply, err := greeter.SayHello(ctx, &pb.HelloRequest{Name: "kei"})
	cancel()
	if err != nil {
		return err
	}
	fmt.Println(reply.Message)

	return nil
}

package main

import (
	"context"
	"os"
	"time"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/grpcflag"
	"github.com/urfave/cli/v2"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var flag = grpcflag.NewServer(
	clix.FlagPrefix("GREETER_"),
	grpcflag.Address("127.0.0.1:9900"),
)

func main() {
	app := cli.NewApp()
	app.Usage = "greeter server"
	app.Flags = flag.Flags()
	app.Before = flag.Before
	app.Action = Action
	exit.Exit(app.Run(os.Args))
}

func Action(c *cli.Context) error {
	s, err := flag.NewServer()
	if err != nil {
		return err
	}
	pb.RegisterGreeterServer(s, &Greeter{})
	return flag.ListenAndServe(s)
}

type Greeter struct {
	pb.UnimplementedGreeterServer
}

func (*Greeter) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	reply := &pb.HelloReply{
		Message: "hello " + req.Name + " ! " + time.Now().Format(time.RFC3339Nano),
	}
	return reply, nil
}

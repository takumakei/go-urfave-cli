Package grpcflag implements functions to use gRPC with github.com/urfave/cli/v2.
======================================================================

[![GoDoc](https://pkg.go.dev/badge/github.com/takumakei/go-urfave-cli/grpcflag)](https://godoc.org/github.com/takumakei/go-urfave-cli/grpcflag)

Server example
----------------------------------------------------------------------

```go
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
```

### output

```
$ go run ./examples/greeter/server --help
NAME:
   server - greeter server

USAGE:
   server [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value, --addr value            address to connect (default: "127.0.0.1:9900") [$GREETER_ADDRESS, $GREETER_ADDR]
   --tls-cert file, --tlscrt file           client certificate pem file [$GREETER_TLS_CERT, $GREETER_TLSCRT]
   --tls-cert-key file, --tlskey file       private key file for client certificate [$GREETER_TLS_CERT_KEY, $GREETER_TLSKEY]
   --tls-gen-cert, --tlsgen                 whether to create and use self signed certificate (default: false) [$GREETER_TLS_GEN_CERT, $GREETER_TLSGEN]
   --tls-verify-client, --mtls              verify client certificate (mTLS) (default: false) [$GREETER_TLS_VERIFY_CLIENT, $GREETER_MTLS]
   --tls-client-ca file, --tlsca file       root CAs certificate file [$GREETER_TLS_CLIENT_CA, $GREETER_TLSCA]
   --tls-min-version value, --tlsmin value  TLS minimum version (default: 1.2) [$GREETER_TLS_MIN_VERSION, $GREETER_TLSMIN]
   --tls-max-version value, --tlsmax value  TLS maximum version (default: 1.3) [$GREETER_TLS_MAX_VERSION, $GREETER_TLSMAX]
   --help, -h                               show help (default: false)
$
```

Client example
----------------------------------------------------------------------

```go
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
```

### output

```
$ go run -trimpath ./examples/greeter/client --help
NAME:
   client - greeter client

USAGE:
   client [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value, --addr value               address to connect (default: "127.0.0.1:9900") [$GREETER_ADDRESS, $GREETER_ADDR]
   --with-insecure, --insecure                 whether to use grpc.WithInsecure (default: false) [$GREETER_WITH_INSECURE, $GREETER_INSECURE]
   --with-block, --block                       whether to use grpc.WithBlock (default: false) [$GREETER_WITH_BLOCK, $GREETER_BLOCK]
   --tls-root-ca file, --tlsca file            root CAs certificate file [$GREETER_TLS_ROOT_CA, $GREETER_TLSCA]
   --tls-server-name value, --tlsserver value  ServerName of tls.Config [$GREETER_TLS_SERVER_NAME, $GREETER_TLSSERVER]
   --tls-skip-verify, --tlsinsecure            InsecureSkipVerify of tls.Config (default: false) [$GREETER_TLS_SKIP_VERIFY, $GREETER_TLSINSECURE]
   --tls-cert file, --tlscrt file              client certificate pem file [$GREETER_TLS_CERT, $GREETER_TLSCRT]
   --tls-cert-key file, --tlskey file          private key file for client certificate [$GREETER_TLS_CERT_KEY, $GREETER_TLSKEY]
   --tls-min-version value, --tlsmin value     TLS minimum version (default: 1.2) [$GREETER_TLS_MIN_VERSION, $GREETER_TLSMIN]
   --tls-max-version value, --tlsmax value     TLS maximum version (default: 1.3) [$GREETER_TLS_MAX_VERSION, $GREETER_TLSMAX]
   --help, -h                                  show help (default: false)
$
```

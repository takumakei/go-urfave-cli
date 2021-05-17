Package netflag implements flags for network client and server
======================================================================

[![GoDoc](https://pkg.go.dev/badge/github.com/takumakei/go-urfave-cli/netflag)](https://godoc.org/github.com/takumakei/go-urfave-cli/netflag)

### Server Example

```
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
```

#### output

```
$ go run ./examples/server
NAME:
   server - A new cli application

USAGE:
   server [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value, --addr value            address to listen [$SERVER_ADDRESS, $SERVER_ADDR]
   --tls-cert file, --tlscrt file           certificate file [$SERVER_TLS_CERT, $SERVER_TLSCRT]
   --tls-cert-key file, --tlskey file       private key file of certificate [$SERVER_TLS_CERT_KEY, $SERVER_TLSKEY]
   --tls-gen-cert, --tlsgen                 generate self-signed certificate (default: false) [$SERVER_TLS_GEN_CERT, $SERVER_TLSGEN]
   --tls-ca file, --tlsca file              root CA file for client auth [$SERVER_TLS_CA, $SERVER_TLSCA]
   --tls-min-version value, --tlsmin value  TLS minimum version (default: 1.2) [$SERVER_TLS_MIN_VERSION, $SERVER_TLSMIN]
   --tls-max-version value, --tlsmax value  TLS maximum version (default: 1.3) [$SERVER_TLS_MAX_VERSION, $SERVER_TLSMAX]
   --help, -h                               show help (default: false)
error: Required flag "addr" not set
exit status 1
$
```

### Client Example

```
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
```

#### output

```
$ go run ./examples/client
NAME:
   client - A new cli application

USAGE:
   client [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value, --addr value            address to connect [$CLIENT_ADDRESS, $CLIENT_ADDR]
   --tls-cert file, --tlscrt file           certificate file [$CLIENT_TLS_CERT, $CLIENT_TLSCRT]
   --tls-cert-key file, --tlskey file       private key file of certificate [$CLIENT_TLS_CERT_KEY, $CLIENT_TLSKEY]
   --tls-ca file, --tlsca file              root CA file of server [$CLIENT_TLS_CA, $CLIENT_TLSCA]
   --tls-server-name value, --tlssrv value  server name for verification [$CLIENT_TLS_SERVER_NAME, $CLIENT_TLSSRV]
   --tls-skip-verify, --tlsinsecure         TLS insecure skip verify (default: false) [$CLIENT_TLS_SKIP_VERIFY, $CLIENT_TLSINSECURE]
   --tls-min-version value, --tlsmin value  TLS minimum version (default: 1.2) [$CLIENT_TLS_MIN_VERSION, $CLIENT_TLSMIN]
   --tls-max-version value, --tlsmax value  TLS maximum version (default: 1.3) [$CLIENT_TLS_MAX_VERSION, $CLIENT_TLSMAX]
   --help, -h                               show help (default: false)
error: Required flag "addr" not set
exit status 1
$
```

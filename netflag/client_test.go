package netflag_test

import (
	"fmt"

	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/netflag"
	"github.com/urfave/cli/v2"
)

func ExampleClient() {
	clientFlags := netflag.NewClientName(
		clix.FlagPrefix("EXAMPLE_"), "echo",
		netflag.Network("udp", "*"),
		netflag.Address("127.0.0.1:9000"),
	)

	app := cli.NewApp()
	app.Name = "example"
	app.HelpName = "example"
	app.Flags = clientFlags.Flags()
	app.Before = clientFlags.Before
	app.Action = func(c *cli.Context) error {
		fmt.Println("network[" + clientFlags.Network() + "]")
		fmt.Println("address[" + clientFlags.Address() + "]")
		return nil
	}

	_ = app.Run([]string{"example", "--help"})
	// Output:
	// NAME:
	//    example - A new cli application
	//
	// USAGE:
	//    example [global options] command [command options] [arguments...]
	//
	// COMMANDS:
	//    help, h  Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --echo-network value, --echo-net value             network to connect (default: "udp") [$EXAMPLE_ECHO_NETWORK, $EXAMPLE_ECHO_NET]
	//    --echo-address value, --echo-addr value            address to connect (default: "127.0.0.1:9000") [$EXAMPLE_ECHO_ADDRESS, $EXAMPLE_ECHO_ADDR]
	//    --echo-tls-cert file, --echo-tlscrt file           certificate file [$EXAMPLE_ECHO_TLS_CERT, $EXAMPLE_ECHO_TLSCRT]
	//    --echo-tls-cert-key file, --echo-tlskey file       private key file of certificate [$EXAMPLE_ECHO_TLS_CERT_KEY, $EXAMPLE_ECHO_TLSKEY]
	//    --echo-tls-ca file, --echo-tlsca file              root CA file of server [$EXAMPLE_ECHO_TLS_CA, $EXAMPLE_ECHO_TLSCA]
	//    --echo-tls-server-name value, --echo-tlssrv value  server name for verification [$EXAMPLE_ECHO_TLS_SERVER_NAME, $EXAMPLE_ECHO_TLSSRV]
	//    --echo-tls-skip-verify, --echo-tlsinsecure         TLS insecure skip verify (default: false) [$EXAMPLE_ECHO_TLS_SKIP_VERIFY, $EXAMPLE_ECHO_TLSINSECURE]
	//    --echo-tls-min-version value, --echo-tlsmin value  TLS minimum version (default: 1.2) [$EXAMPLE_ECHO_TLS_MIN_VERSION, $EXAMPLE_ECHO_TLSMIN]
	//    --echo-tls-max-version value, --echo-tlsmax value  TLS maximum version (default: 1.3) [$EXAMPLE_ECHO_TLS_MAX_VERSION, $EXAMPLE_ECHO_TLSMAX]
	//    --help, -h                                         show help (default: false)
}

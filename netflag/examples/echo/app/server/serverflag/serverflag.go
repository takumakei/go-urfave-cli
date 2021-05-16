package serverflag

import (
	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/netflag"
	"github.com/urfave/cli/v2"
)

var (
	FlagPrefix = clix.FlagPrefix("ECHO_SERVER_")

	Flags = clix.Flags(
		FlagListen.Flags(),
		FlagCount,
	)

	FlagListen = netflag.NewServer(FlagPrefix,
		netflag.Network("tcp", "tcp4", "tcp6", "unix"),
		netflag.Address("127.0.0.1:9900"),
		netflag.DisableGenCert,
		netflag.EnableGenCert,
	)

	FlagCount = &cli.IntFlag{
		Name:        "count",
		Aliases:     []string{"c"},
		Usage:       "number of times to repeat sending a message",
		EnvVars:     FlagPrefix.EnvVars("COUNT", "C"),
		FilePath:    FlagPrefix.FilePath("COUNT", "C"),
		Value:       1,
		Destination: new(int),
	}

	FlagSet = clix.NewFlagSet()
)

func Before(c *cli.Context) error {
	delint.Must(FlagSet.Init(c))
	return FlagListen.Before(c)
}

func Count() int {
	return *FlagCount.Destination
}

func LookupCount() (int, bool) {
	return Count(), FlagSet.IsSet(FlagCount)
}

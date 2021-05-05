package flagapp

import (
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

var (
	FlagPrefix = clix.FlagPrefix("FLAGSET_")

	Flags = []cli.Flag{
		FlagCount,
	}

	FlagCount = &cli.IntFlag{
		Name:        "count",
		Aliases:     []string{"c"},
		EnvVars:     FlagPrefix.EnvVars("COUNT"),
		FilePath:    FlagPrefix.FilePath("COUNT"),
		Destination: new(int),
	}

	FlagSet = clix.NewFlagSet()
)

func Count() int {
	return *FlagCount.Destination
}

func LookupCount() (int, bool) {
	return Count(), FlagSet.IsSet(FlagCount)
}

package flagapp

import (
	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

var (
	FlagPrefix = clix.FlagPrefix("FLAGSET_")

	Flags []cli.Flag

	FlagCount *cli.IntFlag

	FlagSet = clix.NewFlagSet()
)

func init() {
	var (
		nameCount = clix.NewFlagName(FlagPrefix, "", "count")
	)

	FlagCount = &cli.IntFlag{
		Name:        nameCount.Name,
		Aliases:     nameCount.Aliases,
		Usage:       "",
		EnvVars:     nameCount.EnvVars,
		FilePath:    nameCount.FilePath,
		Destination: new(int),
	}

	Flags = []cli.Flag{
		FlagCount,
	}
}

func Before(c *cli.Context) error {
	delint.Must(FlagSet.Init(c))
	return nil
}

func After(c *cli.Context) error {
	return nil
}

func Count() int {
	return *FlagCount.Destination
}

func LookupCount() (int, bool) {
	return Count(), FlagSet.IsSet(FlagCount)
}

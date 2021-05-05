package flagsplit

import (
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

var (
	FlagPrefix = clix.FlagPrefix("FLAGSET_SPLIT_")

	Flags = []cli.Flag{
		FlagHorizontal,
		FlagVertical,
	}

	FlagHorizontal = &cli.IntFlag{
		Name:        "horizontal",
		Aliases:     []string{"H"},
		EnvVars:     FlagPrefix.EnvVars("HORIZONTAL"),
		FilePath:    FlagPrefix.FilePath("HORIZONTAL"),
		Destination: new(int),
	}

	FlagVertical = &cli.IntFlag{
		Name:        "vertical",
		Aliases:     []string{"V"},
		EnvVars:     FlagPrefix.EnvVars("VERTICAL"),
		FilePath:    FlagPrefix.FilePath("VERTICAL"),
		Destination: new(int),
	}

	FlagSet = clix.NewFlagSet()

	Direction = FlagSet.NewExclusiveFlags(
		FlagHorizontal,
		FlagVertical,
	)
)

func Horizontal() int {
	return *FlagHorizontal.Destination
}

func Vertical() int {
	return *FlagVertical.Destination
}

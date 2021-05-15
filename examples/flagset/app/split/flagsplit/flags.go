package flagsplit

import (
	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

var (
	FlagPrefix = clix.FlagPrefix("FLAGSET_SPLIT_")

	Flags []cli.Flag

	FlagHorizontal *cli.IntFlag

	FlagVertical *cli.IntFlag

	FlagSet = clix.NewFlagSet()

	Direction *clix.ExclusiveFlags
)

func init() {
	var (
		nameHorizaontal = clix.NewFlagNameAlias(FlagPrefix, "", "horizontal", "H")
		nameVertical    = clix.NewFlagNameAlias(FlagPrefix, "", "vertical", "V")
	)

	FlagHorizontal = &cli.IntFlag{
		Name:        nameHorizaontal.Name,
		Aliases:     nameHorizaontal.Aliases,
		Usage:       "horizontal `width`",
		EnvVars:     nameHorizaontal.EnvVars,
		FilePath:    nameHorizaontal.FilePath,
		Destination: new(int),
	}

	FlagVertical = &cli.IntFlag{
		Name:        nameVertical.Name,
		Aliases:     nameVertical.Aliases,
		Usage:       "vertical `height`",
		EnvVars:     nameVertical.EnvVars,
		FilePath:    nameVertical.FilePath,
		Destination: new(int),
	}

	Flags = []cli.Flag{
		FlagHorizontal,
		FlagVertical,
	}

	Direction = FlagSet.NewExclusiveFlags(
		FlagHorizontal,
		FlagVertical,
	)
}

func Before(c *cli.Context) error {
	delint.Must(FlagSet.Init(c))
	return nil
}

func After(c *cli.Context) error {
	return nil
}

func Horizontal() int {
	return *FlagHorizontal.Destination
}

func Vertical() int {
	return *FlagVertical.Destination
}

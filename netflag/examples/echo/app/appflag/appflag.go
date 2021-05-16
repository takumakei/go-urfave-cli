package appflag

import (
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

var (
	FlagPrefix = clix.FlagPrefix("ECHO_")

	Flags = clix.Flags(
		FlagDebug,
		FlagVersion,
	)

	FlagDebug = &cli.BoolFlag{
		Name:        "debug",
		Aliases:     []string{"verbose", "v"},
		Usage:       "debug mode",
		EnvVars:     FlagPrefix.EnvVars("DEBUG", "VERBOSE", "V"),
		FilePath:    FlagPrefix.FilePath("DEBUG", "VERBOSE", "V"),
		Destination: new(bool),
	}

	FlagVersion = &cli.BoolFlag{
		Name:        "version",
		Aliases:     []string{"V"},
		Usage:       "print the version",
		Destination: new(bool),
	}

	FlagSet = clix.NewFlagSet()
)

func Debug() bool {
	return *FlagDebug.Destination
}

func LookupDebug() (bool, bool) {
	return Debug(), FlagSet.IsSet(FlagDebug)
}

func Version() bool {
	return *FlagVersion.Destination
}

package main

import (
	"fmt"
	"os"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/netflag/examples/echo/app"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	app.App.Version = fmt.Sprintf("%s-%s-%s", version, len7(commit), date)
	exit.Exit(app.App.Run(os.Args))
}

func len7(s string) string {
	if len(s) > 7 {
		return s[:7]
	}
	return s
}

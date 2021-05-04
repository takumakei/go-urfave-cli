Package fish implements EnableishCompletionCommand for github.com/urfave/cli
======================================================================

The command `fish-completion` is out of the box.

Example
----------------------------------------------------------------------

```
package main

import (
	"os"

	"github.com/takumakei/go-urfave-cli/fish"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	fish.EnableFishCompletionCommand(app)
	_ = app.Run(os.Args)
}
```

### fish-completion

```
$ example fish-completion
NAME:
   example fish-completion - managing fish completions

USAGE:
   example fish-completion command [command options]

DESCRIPTION:
   print, install and uninstall the completion script for fish shell

COMMANDS:
   print      print fish completions
   install    install fish completions script
   uninstall  uninstall fish completions script
   help, h    Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

### print

```
$ example fish-completion help print
NAME:
   example fish-completion print - print fish completions

USAGE:
   example fish-completion print [command options]

OPTIONS:
   --no-help, -n         supress the help command and help flags (default: false)
   --prog name, -p name  name for completion filename (default: "example")
```

### install

```
$ example fish-completion help install
NAME:
   example fish-completion install - install fish completions script

USAGE:
   example fish-completion install [command options]

OPTIONS:
   --no-help, -n         supress the help command and help flags (default: false)
   --dir dir, -d dir     fish config dir (default: "$HOME/.config/fish/completions")
   --prog name, -p name  name for completion filename (default: "example")
```

### uninstall

```
$ example fish-completion help uninstall
NAME:
   example fish-completion uninstall - uninstall fish completions script

USAGE:
   example fish-completion uninstall [command options]

OPTIONS:
   --dir dir, -d dir     fish config dir (default: "$HOME/.config/fish/completions")
   --prog name, -p name  name for completion filename (default: "example")
```

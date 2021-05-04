// Package fish implements EnableFishCompletionCommand for github.com/urfave/cli.
package fish

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/renameio"
	"github.com/takumakei/go-exit"
	"github.com/urfave/cli/v2"
)

// EnableFishCompletionCommand append the command to manage the fish
// completion script.
func EnableFishCompletionCommand(app *cli.App) {
	FlagProg.Value = app.Name
	app.Commands = append(app.Commands, FishCompletionCommand)
}

// FishCompletionCommand is the command to manage the fish completion script.
var FishCompletionCommand = &cli.Command{
	Name:        "fish-completion",
	Usage:       "managing fish completions",
	Description: "print, install and uninstall the completion script for fish shell",
	ArgsUsage:   " ",
	Hidden:      true,
	Subcommands: []*cli.Command{
		PrintCommand,
		InstallCommand,
		UninstallCommand,
	},
}

// showFishCompletionCommand is a function that has
// `FishCompletionCommand.Hidden=false`.
// The reason this variable is not a function is to avoid initialization loop.
var showFishCompletionCommand func()

func init() {
	showFishCompletionCommand = func() {
		FishCompletionCommand.Hidden = false
	}
}

// PrintCommand is the subcommand of FishCompletionCommand.
var PrintCommand = &cli.Command{
	Name:      "print",
	Usage:     "print fish completions",
	ArgsUsage: " ",
	Action:    Print,
	Flags: []cli.Flag{
		FlagNoHelp,
		FlagProg,
	},
}

// InstallCommand is the subcommand of FishCompletionCommand.
var InstallCommand = &cli.Command{
	Name:      "install",
	Usage:     "install fish completions script",
	ArgsUsage: " ",
	Action:    Install,
	Flags: []cli.Flag{
		FlagNoHelp,
		FlagDir,
		FlagProg,
	},
}

// UninstallCommand is the subcommand of FishCompletionCommand.
var UninstallCommand = &cli.Command{
	Name:      "uninstall",
	Usage:     "uninstall fish completions script",
	ArgsUsage: " ",
	Action:    Uninstall,
	Flags: []cli.Flag{
		FlagDir,
		FlagProg,
	},
}

// FlagDir is the flag of the directory to install or uninstall the completion
// script.
var FlagDir = &cli.StringFlag{
	Name:    "dir",
	Usage:   "fish config `dir`",
	Aliases: []string{"d"},
	Value:   "$HOME/.config/fish/completions",
}

// FlagProg is the flag of the base name of the completion script to install or
// uninstall.
var FlagProg = &cli.StringFlag{
	Name:    "prog",
	Usage:   "`name` for completion filename",
	Aliases: []string{"p"},
}

// FlagNoHelp is the flag to suppress the help command and help flags in the
// completion script.
var FlagNoHelp = &cli.BoolFlag{
	Name:        "no-help",
	Usage:       "supress the help command and help flags",
	Aliases:     []string{"n"},
	Destination: new(bool),
}

// Dir returns the value of FlagDir, the directory to install or uninstall the
// completion script.
func Dir(c *cli.Context) (string, error) {
	v := c.String(FlagDir.Name)
	if len(v) == 0 {
		return "", errors.New("empty string for dir")
	}
	return os.ExpandEnv(v), nil
}

// Prog returns the value of FlagProg, the base name of the completion script
// to install or uninstall.
func Prog(c *cli.Context) (string, error) {
	v := c.String(FlagProg.Name)
	if len(v) == 0 {
		v = c.App.Name
		if len(v) == 0 {
			return "", errors.New("empty string for prog")
		}
	}
	return v, nil
}

// ScriptPathname returns the pathname of the completion script, equivalent to
// `filepath.Join(Dir(), Prog() + ".fish")`.
func ScriptPathname(c *cli.Context) (string, error) {
	var path, dir, prog string
	var err error
	if dir, err = Dir(c); err == nil {
		if prog, err = Prog(c); err == nil {
			path = filepath.Join(dir, prog+".fish")
		}
	}
	return path, err
}

// NoHelp returns the value of FlagNoHelp.
func NoHelp() bool {
	return *FlagNoHelp.Destination
}

// Print is the action of PrintCommand, prints fish completions.
func Print(c *cli.Context) error {
	if c.Args().Present() {
		_ = cli.ShowSubcommandHelp(c)
		return exit.Status(1)
	}

	s, err := ToFishCompletion(c, NoHelp())
	if err == nil {
		_, err = c.App.Writer.Write([]byte(s))
	}
	return err
}

// ToFishCompletion creates a fish completion string for the `*App` of the root
// context from c. The function errors if either parsing or writing of the
// string fails.
func ToFishCompletion(c *cli.Context, hideHelp bool) (string, error) {
	var app *cli.App
	for _, v := range c.Lineage() {
		if v.App != nil {
			app = v.App
		}
	}

	prog, err := Prog(c)
	if err != nil {
		return "", err
	}
	app.Name = prog

	// to include the command in the completion script.
	showFishCompletionCommand()

	if hideHelp {
		// to remove the help command and help flags.
		removeHelp(app)
	}

	return app.ToFishCompletion()
}

func removeHelp(app *cli.App) {
	app.HideHelp = true
	app.HideHelpCommand = true
	app.Commands = removeHelpCommand(app.Commands)
	app.Flags = removeHelpFlag(app.Flags)
}

func removeHelpCommand(list []*cli.Command) []*cli.Command {
	cmds := make([]*cli.Command, 0, len(list))
	for _, v := range list {
		if v.Name != "help" {
			v.HideHelp = true
			v.HideHelpCommand = true
			v.Flags = removeHelpFlag(v.Flags)
			v.Subcommands = removeHelpCommand(v.Subcommands)
			cmds = append(cmds, v)
		}
	}
	return cmds
}

func removeHelpFlag(list []cli.Flag) []cli.Flag {
	for i, v := range list {
		if contains("help", v.Names()) {
			a := make([]cli.Flag, len(list)-1)
			copy(a, list[:i])
			copy(a[i:], list[i+1:])
			return a
		}
	}
	return list
}

func contains(s string, a []string) bool {
	for _, v := range a {
		if s == v {
			return true
		}
	}
	return false
}

// Install is the action of InstallCommand, installs the fish completion
// script.
func Install(c *cli.Context) error {
	if c.Args().Present() {
		_ = cli.ShowSubcommandHelp(c)
		return exit.Status(1)
	}

	s, err := ToFishCompletion(c, NoHelp())
	if err != nil {
		return err
	}
	path, err := ScriptPathname(c)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "install %s\n", path)
	return renameio.WriteFile(path, []byte(s), 0644)
}

// Uninstall is the action of UninstallCommand, removes the fish completion
// script.
func Uninstall(c *cli.Context) error {
	if c.Args().Present() {
		_ = cli.ShowSubcommandHelp(c)
		return exit.Status(1)
	}

	path, err := ScriptPathname(c)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "uninstall %s\n", path)
	return os.Remove(path)
}

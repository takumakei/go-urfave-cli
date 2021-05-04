package fish_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takumakei/go-urfave-cli/fish"
	"github.com/urfave/cli/v2"
)

func Example() {
	var app *cli.App = newApp()
	fish.EnableFishCompletionCommand(app)
	_ = app.Run([]string{"example", "fish-completion"})
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "example"
	app.HelpName = "example"
	return app
}

func TestEnableFishCompletionCommand(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		want := readTestData("print.txt")

		s := new(strings.Builder)
		app := cli.NewApp()
		app.Name = "example"
		app.HelpName = "example"
		app.Writer = s
		app.ErrWriter = s
		fish.EnableFishCompletionCommand(app)
		_ = app.Run([]string{"example", "fish-completion", "print"})

		got := s.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("-want +got\n%s", diff)
		}
	})

	t.Run("NoHelp", func(t *testing.T) {
		want := readTestData("print.nohelp.txt")

		s := new(strings.Builder)
		app := cli.NewApp()
		app.Name = "example"
		app.HelpName = "example"
		app.Writer = s
		app.ErrWriter = s
		fish.EnableFishCompletionCommand(app)
		_ = app.Run([]string{"example", "fish-completion", "print", "--no-help"})

		got := s.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("-want +got\n%s", diff)
		}
	})
}

func readTestData(name string) string {
	p, err := ioutil.ReadFile(filepath.Join("test", name))
	if err != nil {
		panic(err)
	}
	return string(p)
}

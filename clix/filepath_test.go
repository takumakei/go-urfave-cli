package clix_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

func ExampleFilePath() {
	// This flag can be set in 3 ways.
	//
	//   1. a command line argument `--password`
	//   2. an environment variable `EXAMPLE_PASSWORD`
	//   3. one of the first file that exists
	//     1. pathname by the environment variable `EXAMPLE_PASSWORD_FILE`
	//     2. or pathname of `$HOME/.config/example/defaults/password`
	//     3. or pathname of `/etc/example/defaults/password`
	flag := &cli.StringFlag{
		Name:    "password",
		EnvVars: []string{"EXAMPLE_PASSWORD"},
		FilePath: clix.FilePath(
			os.Getenv("EXAMPLE_PASSWORD_FILE"),
			os.ExpandEnv("$HOME/.config/example/defaults/password"),
			"/etc/example/defaults/password",
		),
	}
	_ = flag
}

func TestFilePath(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	f1st := filepath.Join(dir, "1st")
	f2nd := filepath.Join(dir, "2nd")
	f3rd := filepath.Join(dir, "3rd")

	t.Run("Default", func(t *testing.T) {
		// no file, no env, no arg
		want := "default"
		testFilePath(t, want, nil, f2nd, f3rd)
	})

	t.Cleanup(func() { os.Unsetenv("EXAMPLE_PASSWORD_FILE") })
	t.Run("ByFilePath", func(t *testing.T) {
		t.Run("3rd", func(t *testing.T) {
			// only 3rd file exists
			want := "3rd"
			if err := ioutil.WriteFile(f3rd, []byte(want), 0644); err != nil {
				t.Fatal(err)
			}
			testFilePath(t, want, nil, f2nd, f3rd)
		})
		t.Run("2nd", func(t *testing.T) {
			// 2nd and 3rd files exist
			want := "2nd"
			if err := ioutil.WriteFile(f2nd, []byte(want), 0644); err != nil {
				t.Fatal(err)
			}
			testFilePath(t, want, nil, f2nd, f3rd)
		})
		t.Run("1st", func(t *testing.T) {
			// 1st, 2nd and 3rd files exist
			want := "1st"
			if err := ioutil.WriteFile(f1st, []byte(want), 0644); err != nil {
				t.Fatal(err)
			}
			os.Setenv("EXAMPLE_PASSWORD_FILE", f1st)
			testFilePath(t, want, nil, f2nd, f3rd)
		})
	})

	t.Cleanup(func() { os.Unsetenv("EXAMPLE_PASSWORD") })
	t.Run("ByEnv", func(t *testing.T) {
		want := "environment variable"
		os.Setenv("EXAMPLE_PASSWORD", want)
		testFilePath(t, want, nil, f2nd, f3rd)
	})
	t.Run("ByArg", func(t *testing.T) {
		want := "command line"
		testFilePath(t, want, []string{"--password", want}, f2nd, f3rd)
	})
}

func testFilePath(t *testing.T, want string, args []string, f2nd, f3rd string) {
	t.Helper()

	got := ""
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "password",
			EnvVars: []string{"EXAMPLE_PASSWORD"},
			FilePath: clix.FilePath(
				os.Getenv("EXAMPLE_PASSWORD_FILE"),
				f2nd,
				f3rd,
			),
			Value: "default",
		},
	}
	app.Action = func(c *cli.Context) error {
		got = c.String("password")
		return nil
	}
	_ = app.Run(append([]string{"prog"}, args...))
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

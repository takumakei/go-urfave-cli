package clix_test

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

func TestFlagSet_Exclusive(t *testing.T) {
	t.Run("withoutEnv", func(t *testing.T) {
		tt := []struct {
			Args []string
			Want string
			Err  string
		}{
			{[]string{"prog"}, "no flag", ""},
			{[]string{"prog", "-a", "1"}, "1", ""},
			{[]string{"prog", "-b", "2"}, "  2", ""},
			{[]string{"prog", "-a", "1", "-b", "2"}, "", "more than one flags are set in args (-a,-b)"},
		}
		for i, v := range tt {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				testFlagSetExclusive(t, v.Args, v.Want, v.Err)
			})
		}
	})

	t.Run("withEnv", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			tt := []struct {
				Args []string
				Want string
				Err  string
			}{
				{[]string{"prog"}, "3", ""},
				{[]string{"prog", "-a", "1"}, "1", ""},
				{[]string{"prog", "-b", "2"}, "  2", ""},
				{[]string{"prog", "-a", "1", "-b", "2"}, "", "more than one flags are set in args (-a,-b)"},
			}
			os.Setenv("APP_A", "3")
			t.Cleanup(func() { os.Unsetenv("APP_A") })
			for i, v := range tt {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					testFlagSetExclusive(t, v.Args, v.Want, v.Err)
				})
			}
		})

		t.Run("2", func(t *testing.T) {
			tt := []struct {
				Args []string
				Want string
				Err  string
			}{
				{[]string{"prog"}, "3", "more than one flags are set in envs (-a,-b)"},
				{[]string{"prog", "-a", "1"}, "1", ""},
				{[]string{"prog", "-b", "2"}, "  2", ""},
				{[]string{"prog", "-a", "1", "-b", "2"}, "", "more than one flags are set in args (-a,-b)"},
			}
			os.Setenv("APP_A", "3")
			os.Setenv("APP_B", "4")
			t.Cleanup(func() {
				os.Unsetenv("APP_A")
				os.Unsetenv("APP_B")
			})
			for i, v := range tt {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					testFlagSetExclusive(t, v.Args, v.Want, v.Err)
				})
			}
		})
	})
}

func testFlagSetExclusive(t *testing.T, args []string, want, wantErr string) {
	t.Helper()

	flagA := &cli.IntFlag{
		Name:    "a",
		EnvVars: []string{"APP_A"},
		FilePath: clix.FilePath(
			os.Getenv("APP_A_FILE"),
		),
	}

	flagB := &cli.IntFlag{
		Name:    "b",
		EnvVars: []string{"APP_B"},
		FilePath: clix.FilePath(
			os.Getenv("APP_B_FILE"),
		),
	}

	flagSet := clix.NewFlagSet()

	app := cli.NewApp()
	app.Flags = []cli.Flag{flagA, flagB}
	var s strings.Builder
	app.Writer = &s
	app.ErrWriter = &s
	app.Before = flagSet.Init
	app.Action = func(c *cli.Context) error {
		flag, err := flagSet.Exclusive(flagA, flagB)
		switch flag {
		case flagA:
			fmt.Fprintf(c.App.Writer, "%d", c.Int("a"))
		case flagB:
			fmt.Fprintf(c.App.Writer, "  %d", c.Int("b"))
		default:
			fmt.Fprintf(c.App.Writer, "no flag")
		}
		return err
	}
	err := app.Run(args)
	if err != nil {
		if diff := cmp.Diff(wantErr, err.Error()); diff != "" {
			t.Fatalf("-want +got\n%s", diff)
		}
		return
	}
	if len(wantErr) > 0 {
		t.Errorf("got no error, want error %q", wantErr)
	}
	got := s.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("-want +got\n%s", diff)
	}
}

func TestErrExclusiveFlags(t *testing.T) {
	errExcl := clix.ErrExclusiveFlags
	errArgs := clix.ErrExclusiveFlagsInArgs
	errEnvs := clix.ErrExclusiveFlagsInEnvs

	t.Run("A", func(t *testing.T) {
		want := true
		got := errors.Is(errArgs, errExcl)
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("E", func(t *testing.T) {
		want := true
		got := errors.Is(errEnvs, errExcl)
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

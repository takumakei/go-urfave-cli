package clix_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

func ExampleChain() {
	flagSet := clix.NewFlagSet()
	app := cli.NewApp()
	app.Before = clix.Chain(flagSet.Init, init1)
}

func init1(c *cli.Context) error {
	return nil
}

func TestChain(t *testing.T) {
	want := []int{0, 1}

	var got []int
	a0 := func(*cli.Context) error {
		got = append(got, 0)
		return nil
	}
	a1 := func(*cli.Context) error {
		got = append(got, 1)
		return exit.Status(42) // abort chain
	}
	a2 := func(*cli.Context) error {
		got = append(got, 2)
		return nil
	}

	f := clix.Chain(a0, a1, a2)
	f(nil)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("-want +got\n%s", diff)
	}
}

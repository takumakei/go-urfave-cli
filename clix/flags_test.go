package clix_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

func TestFlags(t *testing.T) {
	var (
		a = &cli.IntFlag{Name: "a"}
		b = &cli.IntFlag{Name: "b"}
		c = &cli.BoolFlag{Name: "c"}
	)

	want := []cli.Flag{a, b, c}
	got := clix.Flags(a, []cli.Flag{b, c})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("-want +got\n%s", diff)
	}
}

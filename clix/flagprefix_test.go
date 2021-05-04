package clix_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takumakei/go-urfave-cli/clix"
)

func TestFlagPrefix(t *testing.T) {
	flagPrefix := clix.FlagPrefix("HELLO_")

	t.Run("EnvVars", func(t *testing.T) {
		want := []string{"HELLO_FOO"}
		got := flagPrefix.EnvVars("FOO")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("-want +got\n%s", diff)
		}
	})

	t.Run("FilePath", func(t *testing.T) {
		want := os.Args[0]
		os.Setenv("HELLO_FOO_FILE", want)
		t.Cleanup(func() { os.Unsetenv("HELLO_FOO_FILE") })
		got := flagPrefix.FilePath("FOO")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("-want +got\n%s", diff)
		}
	})
}

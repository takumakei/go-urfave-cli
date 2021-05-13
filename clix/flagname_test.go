package clix_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takumakei/go-urfave-cli/clix"
)

func TestNewFlagName(t *testing.T) {
	cases := []struct {
		Prefix clix.FlagPrefix
		Group  string
		Name   string
		File   string
		Want   *clix.FlagName
	}{
		{
			clix.FlagPrefix("EXAMPLE_"),
			"hello",
			"flag",
			"EXAMPLE_HELLO_FLAG_FILE",
			&clix.FlagName{
				Name:     "hello-flag",
				Aliases:  []string{"hello-f"},
				EnvVars:  []string{"EXAMPLE_HELLO_FLAG", "EXAMPLE_HELLO_F"},
				FilePath: os.Args[0],
			},
		},

		{
			clix.FlagPrefix("EXAMPLE_"),
			"hello",
			"flag",
			"EXAMPLE_HELLO_F_FILE",
			&clix.FlagName{
				Name:     "hello-flag",
				Aliases:  []string{"hello-f"},
				EnvVars:  []string{"EXAMPLE_HELLO_FLAG", "EXAMPLE_HELLO_F"},
				FilePath: os.Args[0],
			},
		},

		{
			clix.FlagPrefix(""),
			"hello",
			"flag",
			"HELLO_FLAG_FILE",
			&clix.FlagName{
				Name:     "hello-flag",
				Aliases:  []string{"hello-f"},
				EnvVars:  []string{"HELLO_FLAG", "HELLO_F"},
				FilePath: os.Args[0],
			},
		},

		{
			clix.FlagPrefix(""),
			"",
			"flag",
			"FLAG_FILE",
			&clix.FlagName{
				Name:     "flag",
				Aliases:  []string{"f"},
				EnvVars:  []string{"FLAG", "F"},
				FilePath: os.Args[0],
			},
		},
	}

	file := os.Args[0]
	for i, c := range cases {
		c := c
		os.Setenv(c.File, file)
		t.Cleanup(func() { os.Unsetenv(c.File) })
		got := clix.NewFlagName(c.Prefix, c.Group, c.Name)
		if diff := cmp.Diff(c.Want, got); diff != "" {
			t.Errorf("[%d] -want +got\n%s", i, diff)
		}
	}
}

func TestNewFlagNameAlias(t *testing.T) {
	cases := []struct {
		Prefix clix.FlagPrefix
		Group  string
		Name   string
		Alias  string
		File   string
		Want   *clix.FlagName
	}{
		{
			clix.FlagPrefix("EXAMPLE_"),
			"hello",
			"flag",
			"f-g",
			"EXAMPLE_HELLO_FLAG_FILE",
			&clix.FlagName{
				Name:     "hello-flag",
				Aliases:  []string{"hello-f-g"},
				EnvVars:  []string{"EXAMPLE_HELLO_FLAG", "EXAMPLE_HELLO_F_G"},
				FilePath: os.Args[0],
			},
		},

		{
			clix.FlagPrefix("EXAMPLE_"),
			"hello",
			"flag",
			"f-g",
			"EXAMPLE_HELLO_F_FILE",
			&clix.FlagName{
				Name:     "hello-flag",
				Aliases:  []string{"hello-f-g"},
				EnvVars:  []string{"EXAMPLE_HELLO_FLAG", "EXAMPLE_HELLO_F_G"},
				FilePath: os.Args[0],
			},
		},

		{
			clix.FlagPrefix(""),
			"hello",
			"flag",
			"f-g",
			"HELLO_FLAG_FILE",
			&clix.FlagName{
				Name:     "hello-flag",
				Aliases:  []string{"hello-f-g"},
				EnvVars:  []string{"HELLO_FLAG", "HELLO_F_G"},
				FilePath: os.Args[0],
			},
		},

		{
			clix.FlagPrefix(""),
			"",
			"flag",
			"f-g",
			"FLAG_FILE",
			&clix.FlagName{
				Name:     "flag",
				Aliases:  []string{"f-g"},
				EnvVars:  []string{"FLAG", "F_G"},
				FilePath: os.Args[0],
			},
		},
	}

	file := os.Args[0]
	for i, c := range cases {
		c := c
		os.Setenv(c.File, file)
		t.Cleanup(func() { os.Unsetenv(c.File) })
		got := clix.NewFlagNameAlias(c.Prefix, c.Group, c.Name, c.Alias)
		if diff := cmp.Diff(c.Want, got); diff != "" {
			t.Errorf("[%d] -want +got\n%s", i, diff)
		}
	}
}

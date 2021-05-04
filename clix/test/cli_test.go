package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
)

func Test_flag_1(t *testing.T) {
	if ok, err := strconv.ParseBool(os.Getenv("CLIX_STUDY_CLI")); err != nil || !ok {
		t.SkipNow()
	}

	const (
		T = true
		F = false
	)
	tt := []struct {
		Env  string
		File string
		Arg  string
		L    bool   // c.LocalFlagNames
		F    bool   // flag.IsSet
		C    bool   // c.IsSet
		S    string // result
	}{
		{"E0", "F0", "A0", T, T, T, "A0"},
		{"E1", "F1", "  ", F, T, T, "E1"},
		{"E2", "  ", "A2", T, T, T, "A2"},
		{"E3", "  ", "  ", F, T, T, "E3"},
		{"  ", "F4", "A4", T, T, T, "A4"},
		{"  ", "F5", "  ", F, T, T, "F5"},
		{"  ", "  ", "A6", T, F, T, "A6"},
		{"  ", "  ", "  ", F, F, F, "D7"},
	}

	for i, v := range tt {
		flag := &cli.StringFlag{
			Name:        "X",
			EnvVars:     []string{"FLAG_X"},
			FilePath:    filepath.Join(t.TempDir(), v.File),
			Value:       "D" + strconv.Itoa(i),
			Destination: new(string),
		}
		var got struct {
			L bool
			F bool
			C bool
			S string
		}
		app := &cli.App{
			Flags: []cli.Flag{flag},
			Before: func(c *cli.Context) error {
				for _, s := range c.LocalFlagNames() {
					if s == flag.Name {
						got.L = true
						break
					}
				}
				got.F = flag.IsSet()
				got.C = c.IsSet(flag.Name)
				got.S = c.String(flag.Name)
				return nil
			},
			Action: func(c *cli.Context) error { return nil },
		}
		args := []string{"prog"}
		if s := strings.TrimSpace(v.Env); len(s) > 0 {
			os.Setenv("FLAG_X", s)
		} else {
			os.Unsetenv("FLAG_X")
		}
		if s := strings.TrimSpace(v.File); len(s) > 0 {
			_ = ioutil.WriteFile(flag.FilePath, []byte("F"+strconv.Itoa(i)), 0644)
		}
		if s := strings.TrimSpace(v.Arg); len(s) > 0 {
			args = append(args, "-X", v.Arg)
		}
		_ = app.Run(args)

		if got.L != v.L {
			t.Errorf("%d flag.LocalFlagNames=%v, want %v", i, got.L, v.L)
		}
		if got.F != v.F {
			t.Errorf("%d flag.IsSet=%v, want %v", i, got.F, v.F)
		}
		if got.C != v.C {
			t.Errorf("%d c.IsSet=%v, want %v", i, got.C, v.C)
		}
		if got.S != v.S {
			t.Errorf("%d got=%v, want=%v", i, got.S, v.S)
		}
		if *flag.Destination != v.S {
			t.Errorf("%d Destination=%v, want=%v", i, *flag.Destination, v.S)
		}
		if testing.Verbose() {
			t.Log(i, v, got, *flag.Destination)
		}
	}
}

func Test_flag_2(t *testing.T) {
	if ok, err := strconv.ParseBool(os.Getenv("CLIX_STUDY_CLI")); err != nil || !ok {
		t.SkipNow()
	}

	s := new(strings.Builder)

	helloFlagU := &cli.StringFlag{
		Name:        "u",
		Value:       "DH",
		EnvVars:     []string{"APP_HELLO_U"},
		Destination: new(string),
	}
	hello := &cli.Command{
		Name: "hello",
		Before: func(c *cli.Context) error {
			fmt.Fprintln(s, "hello.Before", c.String("u"), *helloFlagU.Destination)
			return nil
		},
		Action: func(c *cli.Context) error {
			return nil
		},
		Flags: []cli.Flag{helloFlagU},
	}
	appFlagU := &cli.StringFlag{
		Name:        "u",
		Value:       "DA",
		EnvVars:     []string{"APP_U"},
		Destination: new(string),
	}
	app := cli.NewApp()
	app.Flags = []cli.Flag{appFlagU}
	app.Before = func(c *cli.Context) error {
		fmt.Fprintln(s, "app.Before", c.String("u"), *appFlagU.Destination)
		return nil
	}
	app.Commands = []*cli.Command{
		hello,
	}
	os.Setenv("APP_HELLO_U", "E")
	_ = app.Run([]string{
		"prog",
		//"-u", "alice",
		"hello",
		//"-u", "bob",
	})

	want := "app.Before DA DA\nhello.Before E E\n"
	got := s.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("-want +got\n%s", diff)
	}
}

func Test_flag_3(t *testing.T) {
	if ok, err := strconv.ParseBool(os.Getenv("CLIX_STUDY_CLI")); err != nil || !ok {
		t.SkipNow()
	}

	s := new(strings.Builder)

	flagU := &cli.StringSliceFlag{
		Name:    "u",
		EnvVars: []string{"APP_U"},
	}
	hello := &cli.Command{
		Name: "hello",
		Before: func(c *cli.Context) error {
			fmt.Fprintln(s, "hello.Before", strings.Join(c.StringSlice("u"), ":"))
			return nil
		},
		Action: func(c *cli.Context) error { return nil },
		Flags:  []cli.Flag{flagU},
	}
	app := cli.NewApp()
	app.Flags = []cli.Flag{flagU}
	app.Commands = []*cli.Command{
		hello,
	}
	os.Setenv("APP_U", "a,b,c")
	_ = app.Run([]string{
		"prog",
		"-u", "A",
		"-u", "B",
		"hello",
		"-u", "C",
	})

	want := "hello.Before C\n"
	got := s.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("-want +got\n%s", diff)
	}
}

func Test_flag_4(t *testing.T) {
	if ok, err := strconv.ParseBool(os.Getenv("CLIX_STUDY_CLI")); err != nil || !ok {
		t.SkipNow()
	}

	s := new(strings.Builder)

	flagU := &cli.StringSliceFlag{
		Name:        "u",
		EnvVars:     []string{"APP_U"},
		Destination: new(cli.StringSlice),
	}
	helloFlagU := &cli.StringSliceFlag{
		Name:        "u",
		EnvVars:     []string{"APP_HELLO_U"},
		Destination: new(cli.StringSlice),
	}
	hello := &cli.Command{
		Name: "hello",
		Before: func(c *cli.Context) error {
			fmt.Fprintln(s, "hello.Before", strings.Join(flagU.Destination.Value(), ":"))
			fmt.Fprintln(s, "hello.Before", strings.Join(helloFlagU.Destination.Value(), ":"))
			return nil
		},
		Action: func(c *cli.Context) error { return nil },
		Flags:  []cli.Flag{helloFlagU},
	}
	app := cli.NewApp()
	app.Flags = []cli.Flag{flagU}
	app.Commands = []*cli.Command{
		hello,
	}
	os.Setenv("APP_U", "a,b,c")
	os.Setenv("APP_HELLO_U", "D,E,F")
	_ = app.Run([]string{
		"prog",
		//"-u", "A",
		//"-u", "B",
		"hello",
		//"-u", "C",
	})

	want := "hello.Before a:b:c\nhello.Before D:E:F\n"
	got := s.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("-want +got\n%s", diff)
	}
}

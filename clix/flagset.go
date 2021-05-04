package clix

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

// FlagSet represents ...
type FlagSet map[string]struct{}

// NewFlagSet returns a FlagSet.
func NewFlagSet() FlagSet {
	return make(FlagSet)
}

// Init initialize fs using c.LocalFlagNames().
func (fs FlagSet) Init(c *cli.Context) error {
	for _, name := range c.LocalFlagNames() {
		fs[name] = struct{}{}
	}
	return nil
}

// IsSetArgs returns true if flag is specified in arguments
// or false otherwise.
func (fs FlagSet) IsSetArgs(flag cli.Flag) bool {
	for _, name := range flag.Names() {
		if _, ok := fs[name]; ok {
			return true
		}
	}
	return false
}

// IsSet returns true if flag is specified or false otherwise.
// `specified` includes following casese.
//   - specified by arguments
//   - specified by EnvVars
//   - specified by FilePath
func (fs FlagSet) IsSet(flag cli.Flag) bool {
	// flag.IsSet() returns true in following cases.
	//   - specified by EnvVars
	//   - specified by FilePath
	// fs.IsSetArgs(flag) returns true if it exists in arguments.
	return flag.IsSet() || fs.IsSetArgs(flag)
}

var (
	// ErrExclusiveFlagsInArgs represents an error
	// where more than one exclusive flags are set
	// at the same time.
	ErrExclusiveFlags = errors.New("more than one flags are set")

	// ErrExclusiveFlagsInArgs represents an error
	// where more than one exclusive flags are set in arguments
	// at the same time.
	ErrExclusiveFlagsInArgs = fmt.Errorf("%w in args", ErrExclusiveFlags)

	// ErrExclusiveFlagsInEnvs represents an error
	// where more than one exclusive flags are set in environment variables
	// at the same time.
	ErrExclusiveFlagsInEnvs = fmt.Errorf("%w in envs", ErrExclusiveFlags)
)

// Exclusive returns the flag only if it is set exclusively or nil if none is set,
// returns err!=nil if and only if multiple flags are set at the same time.
//
// At first command line arguments are searched.
// Next the environment variables (include FilePath) are searched.
// At most one flag can be specified at the same time,
// otherwise err!=nil is returned.
func (fs FlagSet) Exclusive(flag ...cli.Flag) (cli.Flag, error) {
	p := newPicker(len(flag))
	// fs.IsSetArgs(flag) returns true if it exists in arguments.
	n := p.pick(flag, fs.IsSetArgs)
	switch {
	case n == 1:
		f := p.found[0]
		return f, nil
	case n > 1:
		return nil, fmt.Errorf("%w ("+joinNames(p.found)+")", ErrExclusiveFlagsInArgs)
	}

	p.reset()
	// f.IsSet() returns true in following cases.
	//   - specified by EnvVars
	//   - specified by FilePath
	n = p.pick(flag, func(f cli.Flag) bool { return f.IsSet() })
	switch {
	case n == 1:
		f := p.found[0]
		return f, nil
	case n > 1:
		return nil, fmt.Errorf("%w ("+joinNames(p.found)+")", ErrExclusiveFlagsInEnvs)
	}
	return nil, nil
}

type picker struct {
	found []cli.Flag
}

func newPicker(hint int) *picker {
	return &picker{found: make([]cli.Flag, 0, hint)}
}

func (p *picker) reset() {
	p.found = p.found[:0]
}

func (p *picker) pick(flags []cli.Flag, fn func(cli.Flag) bool) int {
	for _, f := range flags {
		if fn(f) {
			p.found = append(p.found, f)
		}
	}
	return len(p.found)
}

func joinNames(flags []cli.Flag) string {
	names := make([]string, len(flags))
	for i, f := range flags {
		names[i] = "-" + f.Names()[0]
	}
	return strings.Join(names, ",")
}

// Select calls the callback function for exclusive flag that is set,
// returns the result of the callback.
func (fs FlagSet) Select(flags []cli.Flag, cb map[cli.Flag]func() error) error {
	flag, err := fs.Exclusive(flags...)
	if err == nil {
		if fn := cb[flag]; fn != nil {
			err = fn()
		}
	}
	return err
}

// NewExclusiveFlags returns ExclusiveFlags.
func (fs FlagSet) NewExclusiveFlags(flag ...cli.Flag) *ExclusiveFlags {
	return &ExclusiveFlags{FlagSet: fs, Flags: flag}
}

// ExclusiveFlags is pair of FlagSet and Flags.
type ExclusiveFlags struct {
	FlagSet FlagSet
	Flags   []cli.Flag
}

// Select calls the callback function of the exclusive flag that is set,
// returns the result of the callback.
func (ef ExclusiveFlags) Select(cb map[cli.Flag]func() error) error {
	return ef.FlagSet.Select(ef.Flags, cb)
}

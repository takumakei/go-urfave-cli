package zapflag

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// SyncGlobal calls zap.L().Sync(), returns its result.
// If fn is specified, the result is converted by them before returned.
//
// This is intended to be used as cli.AfterFunc.
// see https://pkg.go.dev/github.com/urfave/cli/v2#AfterFunc
func SyncGlobal(fn ...ErrorFunc) func(*cli.Context) error {
	return func(c *cli.Context) error {
		err := zap.L().Sync()
		for _, f := range fn {
			err = f(err)
		}
		return err
	}
}

// Sync calls Sync method of *zap.Logger of *p, returns its result.
// If fn is specified, the result is converted by them before returned.
//
// This is intended to be used as cli.AfterFunc.
// see https://pkg.go.dev/github.com/urfave/cli/v2#AfterFunc
func Sync(p **zap.Logger, fn ...ErrorFunc) func(*cli.Context) error {
	return func(c *cli.Context) error {
		err := (*p).Sync()
		for _, f := range fn {
			err = f(err)
		}
		return err
	}
}

// ErrorFunc converts err.
type ErrorFunc func(err error) error

// IgnoreError returns always nil.
func IgnoreError(error) error {
	return nil
}

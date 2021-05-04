package clix

import "github.com/urfave/cli/v2"

// Chain returns a function of type `func(*cli.Context) error`
// in which each fn is called in order.
func Chain(fn ...func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		for _, f := range fn {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
}

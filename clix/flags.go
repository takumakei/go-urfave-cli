package clix

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Flags returns []cli.Flag in which flags are flatten,
// each element of flags must be nil or a cli.Flag or []cli.Flag,
// panic if none of them. A slice []cli.Flag can contain nil.
func Flags(flags ...interface{}) []cli.Flag {
	list := make([]cli.Flag, 0, len(flags))
	for i, v := range flags {
		switch v := v.(type) {
		case nil:
		case cli.Flag:
			list = append(list, v)
		case []cli.Flag:
			for _, e := range v {
				if e != nil {
					list = append(list, e)
				}
			}
		default:
			panic(fmt.Sprintf("flags[%d]=%v, not of type cli.Flag nor []cli.Flag", i, v))
		}
	}
	return list
}

// FlagIf returns flag if cond is true, otherwise returns nil.
func FlagIf(cond bool, flag ...cli.Flag) []cli.Flag {
	if cond {
		return flag
	}
	return nil
}

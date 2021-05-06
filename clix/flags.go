package clix

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Flags returns []cli.Flag in which flags are flatten,
// each element of flags must be a cli.Flag or []cli.Flag,
// panic if none of them.
func Flags(flags ...interface{}) []cli.Flag {
	list := make([]cli.Flag, 0, len(flags))
	for i, v := range flags {
		switch v := v.(type) {
		case cli.Flag:
			list = append(list, v)
		case []cli.Flag:
			list = append(list, v...)
		default:
			panic(fmt.Sprintf("flags[%d]=%v, not of type cli.Flag nor []cli.Flag", i, v))
		}
	}
	return list
}

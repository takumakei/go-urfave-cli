package clix

import (
	"os"
)

// FlagPrefix represents the prefix of flags.
type FlagPrefix string

// String returns string(fp).
func (fp FlagPrefix) String() string {
	return string(fp)
}

// EnvVars returns []string{FlagPrefix + s, ...}.
func (fp FlagPrefix) EnvVars(s ...string) []string {
	a := make([]string, len(s))
	for i, v := range s {
		a[i] = string(fp) + v
	}
	return a
}

// FilePath returns clix.FilePath(os.Getenv(FlagPrefix + s + "_FILE"), ...).
func (fp FlagPrefix) FilePath(s ...string) string {
	for _, v := range s {
		v = FilePath(os.Getenv(string(fp) + v + "_FILE"))
		if len(v) > 0 {
			return v
		}
	}
	return ""
}

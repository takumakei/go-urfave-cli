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

// EnvVars returns []string{FlagPrefix + s}.
func (fp FlagPrefix) EnvVars(s string) []string {
	return []string{string(fp) + s}
}

// FilePath returns clix.FilePath(os.Getenv(FlagPrefix + s + "_FILE")).
func (fp FlagPrefix) FilePath(s string) string {
	return FilePath(os.Getenv(string(fp) + s + "_FILE"))
}

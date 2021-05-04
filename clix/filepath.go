package clix

import "os"

// FilePath returns the first file that exists
// or an empty string if no file exists.
func FilePath(file ...string) string {
	for _, v := range file {
		if len(v) > 0 {
			if _, err := os.Stat(v); err == nil {
				return v
			}
		}
	}
	return ""
}

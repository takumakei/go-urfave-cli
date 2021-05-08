package zapflag

import (
	"strings"

	"github.com/takumakei/go-urfave-cli/clix"
)

type names struct {
	name     string
	aliases  []string
	envVars  []string
	filePath string
}

func newNames(name, flag string, prefix clix.FlagPrefix) *names {
	return newNamesAlias(name, flag, shortFor(flag), prefix)
}

func shortFor(name string) string {
	m := strings.FieldsFunc(name, func(r rune) bool { return r == '_' || r == '-' })
	a := make([]rune, len(m))
	for i, v := range m {
		a[i] = []rune(v)[0]
	}
	return string(a)
}

func newNamesAlias(name, flag, alias string, prefix clix.FlagPrefix) *names {
	var flagname string
	if len(name) > 0 {
		flagname = name + "-" + flag
		alias = name + "-" + alias
	} else {
		flagname = flag
	}

	uf := strings.ReplaceAll(strings.ToUpper(flagname), "-", "_")
	ua := strings.ReplaceAll(strings.ToUpper(alias), "-", "_")

	return &names{
		name:     flagname,
		aliases:  []string{alias},
		envVars:  prefix.EnvVars(uf, ua),
		filePath: prefix.FilePath(uf, ua),
	}
}

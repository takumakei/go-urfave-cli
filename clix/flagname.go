package clix

import "strings"

type FlagName struct {
	Name     string
	Aliases  []string
	EnvVars  []string
	FilePath string
}

// NewFlagName returns *FlagName.
func NewFlagName(prefix FlagPrefix, group, name string) *FlagName {
	return NewFlagNameAlias(prefix, group, name, ShortFlagName(name))
}

// ShortFlagName makes a short flag name.
//
//   e.g.
//   "target-id" => "ti
//   "LONG_NAME" => "LN"
func ShortFlagName(name string) string {
	m := strings.FieldsFunc(name, func(r rune) bool { return r == '_' || r == '-' })
	a := make([]rune, len(m))
	for i, v := range m {
		a[i] = []rune(v)[0]
	}
	return string(a)
}

// NewFlagNameAlias returns *FlagName.
func NewFlagNameAlias(prefix FlagPrefix, group, name, alias string) *FlagName {
	if len(group) > 0 {
		name = group + "-" + name
		alias = group + "-" + alias
	}

	uf := strings.ReplaceAll(strings.ToUpper(name), "-", "_")
	ua := strings.ReplaceAll(strings.ToUpper(alias), "-", "_")

	return &FlagName{
		Name:     name,
		Aliases:  []string{alias},
		EnvVars:  prefix.EnvVars(uf, ua),
		FilePath: prefix.FilePath(uf, ua),
	}
}

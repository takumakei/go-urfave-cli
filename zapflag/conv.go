package zapflag

import (
	"net/url"
	"sort"
	"strings"

	"github.com/takumakei/go-zapx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapFields converts v into []zap.Field.
//
// For each key in v, if len(v[key]) is less than 2 the element is a result of zap.String.
// Otherwise the element is a result of zap.Strings.
func ZapFields(v url.Values) []zap.Field {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ff := make([]zap.Field, len(keys))
	for i, key := range keys {
		f := v[key]
		if len(f) > 1 {
			ff[i] = zap.Strings(key, f)
		} else {
			ff[i] = zap.String(key, f[0])
		}
	}
	return ff
}

// ZapTimeEncoder converts v into zapcore.TimeEncoder.
//
//   If v is equal to "rfc3339nano" or "RFC3339Nano", it returns zapcore.RFC3339NanoTimeEncoder.
//   If v is equal to "rfc3339" or "RFC3339", it returns zapcore.RFC3339TimeEncoder.
//   If v is equal to "iso8601" or "ISO8601", it returns zapcore.ISO8601TimeEncoder.
//   If v is equal to "millis", it returns zapcore.EpochMillisTimeEncoder.
//   If v is equal to "nanos", it returns zapcore.EpochNanosTimeEncoder.
//   If v is equal to "unix", "epoch" or "", it returns zapcore.EpochTimeEncoder.
//   Otherwise it returns zapcore.TimeEncoderOfLayout(v).
func ZapTimeEncoder(v string) zapcore.TimeEncoder {
	switch v {
	case "rfc3339nano", "RFC3339Nano":
		return zapcore.RFC3339NanoTimeEncoder
	case "rfc3339", "RFC3339":
		return zapcore.RFC3339TimeEncoder
	case "iso8601", "ISO8601":
		return zapcore.ISO8601TimeEncoder
	case "millis":
		return zapcore.EpochMillisTimeEncoder
	case "nanos":
		return zapcore.EpochNanosTimeEncoder
	case "unix", "epoch", "":
		return zapcore.EpochTimeEncoder
	}
	return zapcore.TimeEncoderOfLayout(v)
}

// ParseZapFields converts ss into []zap.Field.
// Each string of ss must have format 'key=value'.
// The key can contain '.'.
func ParseZapFields(ss []string) []zap.Field {
	return newTree(ss).fields()
}

type tree struct {
	m map[string]*tree
	v string
}

func newTree(ss []string) *tree {
	tr := new(tree)
	for _, s := range ss {
		kv := strings.SplitN(s, "=", 2)
		switch len(kv) {
		case 1:
			tr.set(kv[0], "")
		case 2:
			tr.set(kv[0], kv[1])
		}
	}
	return tr
}

func (tr *tree) set(key, value string) {
	m := tr
	for _, k := range strings.Split(key, ".") {
		m = m.getTree(k)
	}
	m.setValue(value)
}

func (tr *tree) getTree(k string) *tree {
	if m, ok := tr.m[k]; ok {
		return m
	}
	if tr.m == nil {
		tr.m = make(map[string]*tree, 1)
		tr.v = ""
	}
	m := new(tree)
	tr.m[k] = m
	return m
}

func (tr *tree) setValue(value string) {
	tr.m = nil
	tr.v = value
}

func (tr *tree) keys() []string {
	ks := make([]string, 0, len(tr.m))
	for k := range tr.m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func (tr *tree) fields() []zap.Field {
	ks := tr.keys()
	zf := make([]zap.Field, len(ks))
	for i, k := range ks {
		m := tr.m[k]
		if m.m == nil {
			zf[i] = zap.String(k, m.v)
		} else {
			zf[i] = zapx.G(k, m.fields()...)
		}
	}
	return zf
}

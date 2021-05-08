package zapflag

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestParseZapFields(t *testing.T) {
	cases := []struct {
		In   []string
		Want string
	}{
		{[]string{"a.a=0", "a.b=1"}, `{"level":"info","msg":"info","a":{"a":"0","b":"1"}}`},
		{[]string{"a.a.a=0", "a.a.b=1"}, `{"level":"info","msg":"info","a":{"a":{"a":"0","b":"1"}}}`},
	}

	for i, c := range cases {
		zf := ParseZapFields(c.In)
		log, s := newLogger()
		log.Info("info", zf...)
		log.Sync()
		got := strings.TrimRightFunc(s.String(), unicode.IsSpace)
		if diff := cmp.Diff(c.Want, got); diff != "" {
			t.Errorf("%d %q\n-want +got\n%s", i, strings.Join(c.In, ","), diff)
		}
	}
}

func newLogger() (*zap.Logger, interface{ String() string }) {
	s := new(strings.Builder)

	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(s),
		zapcore.DebugLevel,
	)

	return zap.New(core), s
}

func TestTree(t *testing.T) {
	cases := []struct {
		In   []string
		Want string
	}{
		{[]string{"a.a=0", "a.b=1"}, `{"a":{"a":"0","b":"1"}}`},
		{[]string{"a.a.a=0", "a.a.b=1"}, `{"a":{"a":{"a":"0","b":"1"}}}`},
	}
	for i, c := range cases {
		tr := newTree(c.In)
		p, err := json.Marshal(tr)
		if err != nil {
			t.Fatal(err)
		}
		got := string(p)
		if diff := cmp.Diff(c.Want, got); diff != "" {
			t.Errorf("%d %q\n-want +got\n%s", i, strings.Join(c.In, ","), diff)
		}
	}
}

func (tr *tree) MarshalJSON() ([]byte, error) {
	if tr.m == nil {
		return json.Marshal(tr.v)
	}
	b := new(bytes.Buffer)
	err := tr.writeJSON(b)
	return b.Bytes(), err
}

func (tr *tree) writeJSON(w *bytes.Buffer) error {
	if tr.m == nil {
		p, err := json.Marshal(tr.v)
		if err == nil {
			_, err = w.Write(p)
		}
		return err
	}

	w.WriteByte('{')
	for i, e := range tr.keys() {
		if i > 0 {
			w.WriteByte(',')
		}
		k, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if _, err := w.Write(k); err != nil {
			return err
		}
		w.WriteByte(':')
		err = tr.m[e].writeJSON(w)
		if err != nil {
			return err
		}
	}
	return w.WriteByte('}')
}

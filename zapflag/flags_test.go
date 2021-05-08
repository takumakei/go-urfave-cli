package zapflag_test

import (
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/zapflag"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_sampling(t *testing.T) {
	cases := []struct {
		Repeat int
		Args   []string
		Want   string
	}{
		{10, []string{"--log-sampling-initial", "3"}, `{"level":"info","msg":"hello","i":0}
{"level":"info","msg":"hello","i":1}
{"level":"info","msg":"hello","i":2}
{"level":"info","msg":"hello","i":5}
{"level":"info","msg":"hello","i":8}
`},

		{10, []string{"--log-sampling-thereafter", "3"}, `{"level":"info","msg":"hello","i":0}
{"level":"info","msg":"hello","i":1}
{"level":"info","msg":"hello","i":2}
{"level":"info","msg":"hello","i":5}
{"level":"info","msg":"hello","i":8}
`},

		{10, []string{"--lsi=3", "--lsth=2"}, `{"level":"info","msg":"hello","i":0}
{"level":"info","msg":"hello","i":1}
{"level":"info","msg":"hello","i":2}
{"level":"info","msg":"hello","i":4}
{"level":"info","msg":"hello","i":6}
{"level":"info","msg":"hello","i":8}
`},

		{10, []string{"--lsi=3", "--lsth=4"}, `{"level":"info","msg":"hello","i":0}
{"level":"info","msg":"hello","i":1}
{"level":"info","msg":"hello","i":2}
{"level":"info","msg":"hello","i":6}
`},
	}
	for i, c := range cases {
		got := sampling(t, c.Repeat, c.Args...)
		if diff := cmp.Diff(c.Want, got); diff != "" {
			t.Log(got)
			t.Errorf("%d %s\n-want +got\n%s", i, strings.Join(c.Args, ","), diff)
		}
	}
}

func sampling(t *testing.T, repeat int, args ...string) string {
	t.Helper()

	s := new(strings.Builder)
	sinkname := "s" + strconv.FormatInt(time.Now().UnixNano(), 16)
	zap.RegisterSink(sinkname, func(*url.URL) (zap.Sink, error) {
		return addSyncClose(s), nil
	})
	os.Setenv("TESTING_LOG_PATH", sinkname+"://")
	os.Setenv("TESTING_LOG_WITH_CALLER", "false")
	os.Setenv("TESTING_LOG_DATE_FORMAT", "")
	t.Cleanup(func() {
		os.Unsetenv("TESTING_LOG_PATH")
		os.Unsetenv("TESTING_LOG_WITH_CALLER")
		os.Unsetenv("TESTING_LOG_DATE_FORMAT")
	})

	zf := zapflag.New(clix.FlagPrefix("TESTING_"))
	app := cli.NewApp()
	app.Flags = zf.Flags()
	app.Before = zf.InitGlobal
	app.Action = func(c *cli.Context) error {
		defer zap.L().Sync()
		for i := 0; i < repeat; i++ {
			zap.L().Info("hello", zap.Int("i", i))
		}
		return nil
	}
	app.Run(append([]string{"testing"}, args...))

	return s.String()
}

func addSyncClose(w io.Writer) zap.Sink {
	return nopWriteCloser{zapcore.AddSync(w)}
}

type nopWriteCloser struct {
	zapcore.WriteSyncer
}

func (nopWriteCloser) Close() error { return nil }

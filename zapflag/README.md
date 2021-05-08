Package zapflag implements functions to use go.uber.org/zap with github.com/urfave/cli
======================================================================

### Example

```
package main

import (
	"os"

	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/zapflag"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	zf := zapflag.New(clix.FlagPrefix("MINI_"))

	app := cli.NewApp()
	app.Flags = zf.Flags()
	app.Before = zf.InitGlobal
	app.After = func(c *cli.Context) error {
		zap.L().Sync()
		return nil
	}
	app.Action = func(c *cli.Context) error {
		zap.L().Info("hello world")
		return nil
	}
	exit.Exit(app.Run(os.Args))
}
```

#### outputs

```
$ mini --help
NAME:
   mini - A new cli application

USAGE:
   mini [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-development, --ld                      enable development mode (default: false) [$MINI_LOG_DEVELOPMENT, $MINI_LD]
   --log-level level, --ll level                level [debug|info|warn|error|dpanic|panic|fatal] (default: auto) [$MINI_LOG_LEVEL, $MINI_LL]
   --log-with-caller, --lwc                     whether including caller (default: auto) [$MINI_LOG_WITH_CALLER, $MINI_LWC]
   --log-stack-trace, --lst                     whether including stack trace (default: auto) [$MINI_LOG_STACK_TRACE, $MINI_LST]
   --log-stack-trace-level level, --lstl level  level [debug|info|warn|error|dpanic|panic|fatal] (default: auto) [$MINI_LOG_STACK_TRACE_LEVEL, $MINI_LSTL]
   --log-date-format value, --ldf value         see https://pkg.go.dev/time#Time.Format (default: auto) [$MINI_LOG_DATE_FORMAT, $MINI_LDF]
   --log-field key=value, --lf key=value        key=value added to the logger [$MINI_LOG_FIELD, $MINI_LF]
   --log-path value, --lp value                 output path (default: auto) [$MINI_LOG_PATH, $MINI_LP]
   --log-err-path value, --lep value            error output path (default: auto) [$MINI_LOG_ERR_PATH, $MINI_LEP]
   --log-sampling-initial N, --lsi N            sampling initial count N (default: auto) [$MINI_LOG_SAMPLING_INITIAL, $MINI_LSI]
   --log-sampling-thereafter N, --lsth N        sampling thereafter count N (default: auto) [$MINI_LOG_SAMPLING_THEREAFTER, $MINI_LSTH]
   --help, -h                                   show help (default: false)
$ mini
{"level":"info","ts":1620542056.145781,"caller":"mini/main.go:24","msg":"hello world"}
$ mini --log-development
2021-05-09T15:34:21.037+0900    INFO    mini/main.go:24 hello world
$ mini --log-development --log-stack-trace-level info
2021-05-09T15:34:30.046+0900    INFO    mini/main.go:24 hello world
main.main.func2
        github.com/takumakei/go-urfave-cli/zapflag/examples/mini/main.go:24
github.com/urfave/cli/v2.(*App).RunContext
        github.com/urfave/cli/v2@v2.3.0/app.go:322
github.com/urfave/cli/v2.(*App).Run
        github.com/urfave/cli/v2@v2.3.0/app.go:224
main.main
        github.com/takumakei/go-urfave-cli/zapflag/examples/mini/main.go:27
runtime.main
        runtime/proc.go:225
$
```

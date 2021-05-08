package zapflag_test

import (
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/takumakei/go-urfave-cli/zapflag"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func Example() {
	zf := zapflag.New(clix.FlagPrefix("EXAMPLE_"))

	app := cli.NewApp()
	app.Name = "example"
	app.HelpName = "example"
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
	app.Run([]string{"example", "--help"})
	// Output:
	// NAME:
	//    example - A new cli application
	//
	// USAGE:
	//    example [global options] command [command options] [arguments...]
	//
	// COMMANDS:
	//    help, h  Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --log-development, --ld                      enable development mode (default: false) [$EXAMPLE_LOG_DEVELOPMENT, $EXAMPLE_LD]
	//    --log-level level, --ll level                level [debug|info|warn|error|dpanic|panic|fatal] (default: auto) [$EXAMPLE_LOG_LEVEL, $EXAMPLE_LL]
	//    --log-with-caller, --lwc                     whether including caller (default: auto) [$EXAMPLE_LOG_WITH_CALLER, $EXAMPLE_LWC]
	//    --log-stack-trace, --lst                     whether including stack trace (default: auto) [$EXAMPLE_LOG_STACK_TRACE, $EXAMPLE_LST]
	//    --log-stack-trace-level level, --lstl level  level [debug|info|warn|error|dpanic|panic|fatal] (default: auto) [$EXAMPLE_LOG_STACK_TRACE_LEVEL, $EXAMPLE_LSTL]
	//    --log-date-format value, --ldf value         see https://pkg.go.dev/time#Time.Format (default: auto) [$EXAMPLE_LOG_DATE_FORMAT, $EXAMPLE_LDF]
	//    --log-field key=value, --lf key=value        key=value added to the logger [$EXAMPLE_LOG_FIELD, $EXAMPLE_LF]
	//    --log-path value, --lp value                 output path (default: auto) [$EXAMPLE_LOG_PATH, $EXAMPLE_LP]
	//    --log-err-path value, --lep value            error output path (default: auto) [$EXAMPLE_LOG_ERR_PATH, $EXAMPLE_LEP]
	//    --log-sampling-initial N, --lsi N            sampling initial count N (default: auto) [$EXAMPLE_LOG_SAMPLING_INITIAL, $EXAMPLE_LSI]
	//    --log-sampling-thereafter N, --lsth N        sampling thereafter count N (default: auto) [$EXAMPLE_LOG_SAMPLING_THEREAFTER, $EXAMPLE_LSTH]
	//    --help, -h                                   show help (default: false)
}

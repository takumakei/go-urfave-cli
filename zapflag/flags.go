package zapflag

import (
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Flags represents flags to build a zap.Logger.
type Flags struct {
	Name string

	FlagDevelopment  *cli.BoolFlag
	FlagLevel        *cli.GenericFlag
	FlagWithCaller   *cli.BoolFlag
	FlagStackTrace   *cli.BoolFlag
	FlagStackTraceLv *cli.GenericFlag
	FlagDateFormat   *cli.StringFlag
	FlagFields       *cli.StringSliceFlag
	FlagPaths        *cli.StringSliceFlag
	FlagErrPaths     *cli.StringSliceFlag

	FlagSamplingInitial    *cli.IntFlag
	FlagSamplingThereafter *cli.IntFlag

	FlagSet clix.FlagSet
}

// New returns NewName("", prefix).
func New(prefix clix.FlagPrefix) *Flags {
	return NewName("", prefix)
}

// NewName returns a *Flags.
func NewName(name string, prefix clix.FlagPrefix) *Flags {
	var (
		logDevelopment        = clix.NewFlagName(prefix, name, "log-development")
		logLevel              = clix.NewFlagName(prefix, name, "log-level")
		logWithCaller         = clix.NewFlagName(prefix, name, "log-with-caller")
		logStackTrace         = clix.NewFlagName(prefix, name, "log-stack-trace")
		logStackTraceLv       = clix.NewFlagName(prefix, name, "log-stack-trace-level")
		logDateFormat         = clix.NewFlagName(prefix, name, "log-date-format")
		logFields             = clix.NewFlagName(prefix, name, "log-field")
		logPaths              = clix.NewFlagName(prefix, name, "log-path")
		logErrPaths           = clix.NewFlagName(prefix, name, "log-err-path")
		logSamplingInitial    = clix.NewFlagName(prefix, name, "log-sampling-initial")
		logSamplingThereafter = clix.NewFlagNameAlias(prefix, name, "log-sampling-thereafter", "lsth")
	)

	return &Flags{
		Name: name,

		FlagDevelopment: &cli.BoolFlag{
			Name:        logDevelopment.Name,
			Aliases:     logDevelopment.Aliases,
			Usage:       "enable development mode",
			EnvVars:     logDevelopment.EnvVars,
			FilePath:    logDevelopment.FilePath,
			Destination: new(bool),
		},

		FlagLevel: &cli.GenericFlag{
			Name:        logLevel.Name,
			Aliases:     logLevel.Aliases,
			Usage:       "`level` [debug|info|warn|error|dpanic|panic|fatal]",
			EnvVars:     logLevel.EnvVars,
			FilePath:    logLevel.FilePath,
			Value:       new(zapcore.Level),
			DefaultText: "auto",
		},

		FlagWithCaller: &cli.BoolFlag{
			Name:        logWithCaller.Name,
			Aliases:     logWithCaller.Aliases,
			Usage:       "whether including caller",
			EnvVars:     logWithCaller.EnvVars,
			FilePath:    logWithCaller.FilePath,
			DefaultText: "auto",
			Destination: new(bool),
		},

		FlagStackTrace: &cli.BoolFlag{
			Name:        logStackTrace.Name,
			Aliases:     logStackTrace.Aliases,
			Usage:       "whether including stack trace",
			EnvVars:     logStackTrace.EnvVars,
			FilePath:    logStackTrace.FilePath,
			DefaultText: "auto",
			Destination: new(bool),
		},

		FlagStackTraceLv: &cli.GenericFlag{
			Name:        logStackTraceLv.Name,
			Aliases:     logStackTraceLv.Aliases,
			Usage:       "`level` [debug|info|warn|error|dpanic|panic|fatal]",
			EnvVars:     logStackTraceLv.EnvVars,
			FilePath:    logStackTraceLv.FilePath,
			Value:       new(zapcore.Level),
			DefaultText: "auto",
		},

		FlagDateFormat: &cli.StringFlag{
			Name:        logDateFormat.Name,
			Aliases:     logDateFormat.Aliases,
			Usage:       "see https://pkg.go.dev/time#Time.Format",
			EnvVars:     logDateFormat.EnvVars,
			FilePath:    logDateFormat.FilePath,
			DefaultText: "auto",
			Destination: new(string),
		},

		FlagFields: &cli.StringSliceFlag{
			Name:        logFields.Name,
			Aliases:     logFields.Aliases,
			Usage:       "`key=value` added to the logger",
			EnvVars:     logFields.EnvVars,
			FilePath:    logFields.FilePath,
			Destination: new(cli.StringSlice),
		},

		FlagPaths: &cli.StringSliceFlag{
			Name:        logPaths.Name,
			Aliases:     logPaths.Aliases,
			Usage:       "output path (default: auto)",
			EnvVars:     logPaths.EnvVars,
			FilePath:    logPaths.FilePath,
			Destination: new(cli.StringSlice),
			//DefaultText: "auto", // FIXME: cli@v2.3.0 does not print DefaultText.
		},

		FlagErrPaths: &cli.StringSliceFlag{
			Name:        logErrPaths.Name,
			Aliases:     logErrPaths.Aliases,
			Usage:       "error output path (default: auto)",
			EnvVars:     logErrPaths.EnvVars,
			FilePath:    logErrPaths.FilePath,
			Destination: new(cli.StringSlice),
			//DefaultText: "auto", // FIXME: cli@v2.3.0 does not print DefaultText.
		},

		FlagSamplingInitial: &cli.IntFlag{
			Name:        logSamplingInitial.Name,
			Aliases:     logSamplingInitial.Aliases,
			Usage:       "sampling initial count `N`",
			EnvVars:     logSamplingInitial.EnvVars,
			FilePath:    logSamplingInitial.FilePath,
			DefaultText: "auto",
			Destination: new(int),
		},

		FlagSamplingThereafter: &cli.IntFlag{
			Name:        logSamplingThereafter.Name,
			Aliases:     logSamplingThereafter.Aliases,
			Usage:       "sampling thereafter count `N`",
			EnvVars:     logSamplingThereafter.EnvVars,
			FilePath:    logSamplingThereafter.FilePath,
			DefaultText: "auto",
			Destination: new(int),
		},

		FlagSet: clix.NewFlagSet(),
	}
}

// Flags returns []cli.Flag.
//
// Elements are following.
//   f.FlagDevelopment
//   f.FlagLevel
//   f.FlagWithCaller
//   f.FlagStackTrace
//   f.FlagStackTraceLv
//   f.FlagDateFormat
//   f.FlagFields
//   f.FlagPaths
//   f.FlagErrPaths
//   f.FlagSamplingInitial
//   f.FlagSamplingThereafter
func (f *Flags) Flags() []cli.Flag {
	return []cli.Flag{
		f.FlagDevelopment,
		f.FlagLevel,
		f.FlagWithCaller,
		f.FlagStackTrace,
		f.FlagStackTraceLv,
		f.FlagDateFormat,
		f.FlagFields,
		f.FlagPaths,
		f.FlagErrPaths,
		f.FlagSamplingInitial,
		f.FlagSamplingThereafter,
	}
}

// InitGlobal calls f.Init(c), then replaces the global logger with f.Logger().
//
// This is intended to be used as cli.BeforeFunc.
// see https://pkg.go.dev/github.com/urfave/cli/v2#BeforeFunc
func (f *Flags) InitGlobal(c *cli.Context) error {
	if err := f.Init(c); err != nil {
		return err
	}

	logger, err := f.Logger()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}

// SyncGlobal calls zap.L().Sync().
//
// This is intended to be used as cli.AfterFunc.
// see https://pkg.go.dev/github.com/urfave/cli/v2#AfterFunc
func SyncGlobal(c *cli.Context) error {
	zap.L().Sync()
	return nil
}

// Init returns f.FlagSet.Init(c).
//
// This is intended to be used as cli.BeforeFunc.
// see https://pkg.go.dev/github.com/urfave/cli/v2#BeforeFunc
func (f *Flags) Init(c *cli.Context) error {
	return f.FlagSet.Init(c)
}

// InitInto returns the function that can be used as cli.BeforeFunc.
// The function calls f.Init and f.Logger to create a new logger.
// Finaly new logger is named with f.Name, saved into p.
func (f *Flags) InitInto(p **zap.Logger) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if err := f.Init(c); err != nil {
			return err
		}

		logger, err := f.Logger()
		if err == nil {
			*p = logger
		}
		return err
	}
}

// Development returns the value of f.FlagDevelopment.
func (f *Flags) Development() bool {
	return *f.FlagDevelopment.Destination
}

// Level returns the value of f.FlagLevel.
func (f *Flags) Level() zapcore.Level {
	return *f.FlagLevel.Value.(*zapcore.Level)
}

// LookupLevel returns the value of f.FlagLevel.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupLevel() (zapcore.Level, bool) {
	return f.Level(), f.FlagSet.IsSet(f.FlagLevel)
}

// WithCaller returns the value of f.FlagWithCaller.
func (f *Flags) WithCaller() bool {
	return *f.FlagWithCaller.Destination
}

// LookupWithCaller returns the value of f.FlagWithCaller.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupWithCaller() (bool, bool) {
	return f.WithCaller(), f.FlagSet.IsSet(f.FlagWithCaller)
}

// StackTrace returns the value of f.FlagStackTrace.
func (f *Flags) StackTrace() bool {
	return *f.FlagStackTrace.Destination
}

// LookupStackTrace returns the value of f.FlagStackTrace.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupStackTrace() (bool, bool) {
	return f.StackTrace(), f.FlagSet.IsSet(f.FlagStackTrace)
}

// StackTraceLv returns the value of f.FlagStackTraceLv.
func (f *Flags) StackTraceLv() zapcore.Level {
	return *f.FlagStackTraceLv.Value.(*zapcore.Level)
}

// LookupStackTraceLv returns the value of f.FlagStackTraceLv.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupStackTraceLv() (zapcore.Level, bool) {
	return f.StackTraceLv(), f.FlagSet.IsSet(f.FlagStackTraceLv)
}

// DateFormat returns the value of f.FlagDateFormat.
func (f *Flags) DateFormat() string {
	return *f.FlagDateFormat.Destination
}

// LookupDateFormat returns the value of f.FlagDateFormat.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupDateFormat() (string, bool) {
	return f.DateFormat(), f.FlagSet.IsSet(f.FlagDateFormat)
}

// Fields returns the value of f.FlagFields.
func (f *Flags) Fields() []string {
	return f.FlagFields.Destination.Value()
}

// LookupFields returns the value of f.FlagFields.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupFields() ([]string, bool) {
	return f.Fields(), f.FlagSet.IsSet(f.FlagFields)
}

// Paths returns the value of f.FlagPaths.
func (f *Flags) Paths() []string {
	return f.FlagPaths.Destination.Value()
}

// LookupPaths returns the value of f.FlagPaths.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupPaths() ([]string, bool) {
	return f.Paths(), f.FlagSet.IsSet(f.FlagPaths)
}

// ErrPaths returns the value of f.FlagErrPaths.
func (f *Flags) ErrPaths() []string {
	return f.FlagErrPaths.Destination.Value()
}

// LookupErrPaths returns the value of f.FlagErrPaths.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupErrPaths() ([]string, bool) {
	return f.ErrPaths(), f.FlagSet.IsSet(f.FlagErrPaths)
}

// SamplingInitial returns the value of f.FlagSamplingInitial.
func (f *Flags) SamplingInitial() int {
	return *f.FlagSamplingInitial.Destination
}

// LookupSamplingInitial returns the value of f.FlagSamplingInitial.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupSamplingInitial() (int, bool) {
	return f.SamplingInitial(), f.FlagSet.IsSet(f.FlagSamplingInitial)
}

// SamplingThereafter returns the value of f.SamplingThereafter.
func (f *Flags) SamplingThereafter() int {
	return *f.FlagSamplingThereafter.Destination
}

// LookupSamplingThereafter returns the value of f.SamplingThereafter.
// If the flag is set by the environment variable or the command line argument the value is returned and the boolean is true.
// Otherwise the returned value is undefined and boolean will be false.
func (f *Flags) LookupSamplingThereafter() (int, bool) {
	return f.SamplingThereafter(), f.FlagSet.IsSet(f.FlagSamplingThereafter)
}

// Logger calls f.RootLogger() and returns the logger named with f.Name if
// f.Name is not empty, otherwise returns the result of f.RootLogger().
func (f *Flags) Logger() (*zap.Logger, error) {
	logger, err := f.RootLogger()
	if len(f.Name) > 0 && logger != nil {
		logger = logger.Named(f.Name)
	}
	return logger, err
}

// RootLogger returns f.Config().Build(f.Options()...).
func (f *Flags) RootLogger() (*zap.Logger, error) {
	return f.Config().Build(f.Options()...)
}

// Config returns zap.Config initialized by f's values.
//
// Following values affects the result.
//   f.Development()
//   f.LookupLevel()
//   f.LookupStackTrace()
//   f.LookupDateFormat()
//   f.LookupPaths()
//   f.LookupErrPaths()
//   f.SamplingConfig()
func (f *Flags) Config() zap.Config {
	config := f.config()

	if v, ok := f.LookupLevel(); ok {
		config.Level = zap.NewAtomicLevelAt(v)
	}

	if v, ok := f.LookupStackTrace(); ok {
		config.DisableStacktrace = !v
	}

	if v, ok := f.LookupDateFormat(); ok {
		if len(v) > 0 {
			config.EncoderConfig.EncodeTime = ZapTimeEncoder(v)
		} else {
			config.EncoderConfig.TimeKey = ""
		}
	}

	if v, ok := f.LookupPaths(); ok {
		config.OutputPaths = v
	}

	if v, ok := f.LookupErrPaths(); ok {
		config.ErrorOutputPaths = v
	}

	if v := f.SamplingConfig(); v != nil {
		config.Sampling = v
	}

	return config
}

func (f *Flags) config() zap.Config {
	if f.Development() {
		return zap.NewDevelopmentConfig()
	}
	return zap.NewProductionConfig()
}

// SamplingConfig returns *zap.SamplingConfig.
// If at least one of f.FlagSamplingInitial or f.FlagSamplingThereafter is set,
// it returns non-nil value.  Otherwise returns nil.
// The value that is not specified is same as the other value.
// The value of thereafter will not be less than 1.
func (f *Flags) SamplingConfig() *zap.SamplingConfig {
	i, iok := f.LookupSamplingInitial()
	t, tok := f.LookupSamplingThereafter()

	if !iok && !tok {
		return nil
	}

	switch {
	case !iok:
		i = t
	case !tok:
		t = i
	}

	if t < 1 {
		t = 1
	}

	return &zap.SamplingConfig{Initial: i, Thereafter: t}
}

// Options returns []zap.Option.
//
// Following values affects the result.
//  - f.LookupStackTraceLv()
//  - f.LookupWithCaller()
//  - f.LookupFields()
func (f *Flags) Options() []zap.Option {
	var opts []zap.Option

	if v, ok := f.LookupStackTraceLv(); ok {
		opts = append(opts, zap.AddStacktrace(v))
	}

	if v, ok := f.LookupWithCaller(); ok {
		opts = append(opts, zap.WithCaller(v))
	}

	if v, ok := f.LookupFields(); ok {
		opts = append(opts, zap.Fields(ParseZapFields(v)...))
	}

	return opts
}

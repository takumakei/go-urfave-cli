package zapflag

import (
	"strings"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// Named adds a new path segment of command/subcommand name to the logger's name.
func Named(logger *zap.Logger, c *cli.Context) *zap.Logger {
	app := strings.ReplaceAll(c.App.Name, " ", ".")
	if len(c.Command.Name) > 0 {
		return logger.Named(app + "." + c.Command.Name)
	}
	return logger.Named(app)
}

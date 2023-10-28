package command

import (
	"os"

	"github.com/zukrin/versifyr/internal/configuration"
	"github.com/zukrin/versifyr/internal/logging"

	"github.com/urfave/cli/v2"
)

// init local configuration
// test if exists
// create if not with a standard configuration with dfferent file types
var InitCommand = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "init project configuration",
	Description: `init creates a new example configuration in the current folder.
It creates .versifyr folder with a configuration.yaml file.`,
	Action:          doInit,
	Flags:           []cli.Flag{},
	HideHelpCommand: true,
	HideHelp:        true,
}

func doInit(cCtx *cli.Context) error {
	cfg := cCtx.App.Metadata["config"].(*configuration.Config)
	logger := cCtx.App.Metadata["logger"].(*logging.Logger)

	// check if exists
	_, err := os.Stat(cfg.BasePath)
	if err == nil {
		logger.Error("Esisting folder %s", cfg.BasePath)
		return err
	}

	// create if not
	err = cfg.CreateConfiguration(logger)
	return err
}

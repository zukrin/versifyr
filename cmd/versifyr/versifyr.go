package main

import (
	"fmt"
	"os"
	"time"

	"github.com/zukrin/versifyr/internal/command"
	"github.com/zukrin/versifyr/internal/configuration"
	"github.com/zukrin/versifyr/internal/logging"
	"github.com/zukrin/versifyr/internal/versifyr"

	"github.com/urfave/cli/v2"
)

// func init() {
// 	cli.VersionPrinter = func(cCtx *cli.Context) {
// 		fmt.Fprintf(cCtx.App.Writer, "version=%s\n", cCtx.App.Version)
// 	}

// 	cli.FlagStringer = func(fl cli.Flag) string {
// 		return fmt.Sprintf("\t\t%s", fl.Names()[0])
// 	}
// }

var cfg *configuration.Config

var logger = logging.NewLogger()

func main() {
	var err error

	cfg, err = configuration.NewConfig()
	if err != nil {
		logger.Info("Warn reading configuration file - %v", err)
	}

	ctm, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", versifyr.Compiled)

	app := &cli.App{
		Name:     "versifyr",
		Version:  versifyr.Version,
		Compiled: ctm,
		Authors: []*cli.Author{
			{
				Name:  "Stefano Zuccaro",
				Email: "zukrin@gmail.com",
			},
		},
		Copyright:       "(c) 2023 Stefano Zuccaro",
		HelpName:        "versifyr",
		HideHelpCommand: false,
		HideHelp:        false,

		Metadata: map[string]interface{}{
			"config": cfg,
			"logger": logger,
		},

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "set output to debug",
				Destination: &cfg.Debug,
				Aliases:     []string{"d"},
				Value:       false,
			},
			&cli.BoolFlag{
				Name:        "nochange",
				Usage:       "simulate changes",
				Destination: &cfg.Simulate,
				Aliases:     []string{"n"},
				DefaultText: "false",
				Value:       false,
			},
		},

		//		Usage:     "",
		//		UsageText: "versifyr - handle versioning into project files",
		//		ArgsUsage: "[args and such]",
		Commands: []*cli.Command{
			command.InitCommand,
			command.ShowCommand,
			command.SetCommand,
		},

		Writer:    logger.InfoWriter,
		ErrWriter: logger.ErrWriter,

		EnableBashCompletion: true,
		HideVersion:          false,

		CommandNotFound: func(cCtx *cli.Context, command string) {
			fmt.Fprintf(cCtx.App.ErrWriter, "Unknown command %q.\n", command)
		},

		ExitErrHandler: func(cCtx *cli.Context, err error) {
			if err != nil {
				fmt.Fprintf(cCtx.App.ErrWriter, "Error: %v\n", err)
				os.Exit(1)
			}
		},

		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(cCtx.App.ErrWriter, "WRONG: %#v\n", err)
			return nil
		},
	}

	app.Run(os.Args)

}

package command

import (
	"bytes"
	"fmt"
	"os"

	"github.com/zukrin/versifyr/internal/configuration"
	"github.com/zukrin/versifyr/internal/logging"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/urfave/cli/v2"
)

// show actual version file by file
var ShowCommand = &cli.Command{
	Name:    "show",
	Aliases: []string{"s"},
	Usage:   "show actual configuration and file content",
	Action:  doShow,
}

func doShow(cCtx *cli.Context) error {

	cfg := cCtx.App.Metadata["config"].(*configuration.Config)
	logger := cCtx.App.Metadata["logger"].(*logging.Logger)

	newlineSW := new(bytes.Buffer)
	newlineSW.WriteString("# actual situation\n")
	for _, file := range cfg.Files {
		newlineSW.WriteString(fmt.Sprintf("## %s\n", file.Name))
		bytes, err := os.ReadFile(file.Path)
		if err != nil {
			newlineSW.WriteString(fmt.Sprintf("> not found, cause %s\n", err.Error()))
		}
		newlineSW.WriteString(fmt.Sprintf("```%s\n", file.Type))
		newlineSW.Write(bytes)
		newlineSW.WriteString("\n```\n")
	}

	result := markdown.Render(newlineSW.String(), 132, 6)
	logger.Info(string(result))

	return nil
}

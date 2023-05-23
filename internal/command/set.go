package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/zukrin/versifyr/internal/configuration"
	"github.com/zukrin/versifyr/internal/logging"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/urfave/cli/v2"

	"github.com/hashicorp/go-version"
)

// change content of files following the configured pattern and using the <k,v> pairs passssas arguments
var SetCommand = &cli.Command{
	Name:    "set",
	Aliases: []string{"s"},
	Usage:   "set values as key=value to be replaced in files",
	Action:  doSet,
}

func doSet(cCtx *cli.Context) error {

	cfg := cCtx.App.Metadata["config"].(*configuration.Config)
	logger := cCtx.App.Metadata["logger"].(*logging.Logger)

	logger.Info("setting values")

	if !cCtx.Args().Present() {
		return errors.New("no values to set")
	}

	dictionary := make(map[string]string)

	// split key and value
	for _, a := range cCtx.Args().Slice() {
		vvv := strings.Split(a, "=")
		if len(vvv) != 2 {
			return errors.New("syntax error defining values in " + a)
		}
		dictionary[vvv[0]] = vvv[1]

	}

	// set default values
	dictionary = setWellKnownValues(dictionary)
	logger.Info("using values %v", dictionary)

	setFiles := make([]*configuration.ConfigFile, 0)

	// replace values in files
	for _, file := range cfg.Files {

		for _, p := range file.Placeholders {
			// for each placeholder replace the designed line with the template output
			newlineSW := new(bytes.Buffer)
			err := p.Template.Execute(newlineSW, dictionary)
			if err != nil {
				return err
			}
			newline := newlineSW.String()
			file.Lines[p.Line] = newline
			logger.Info("replaced into %s line %v with %s", file.Name, p.Line, file.Lines[p.Line])
		}
		setFiles = append(setFiles, file)
	}

	// write back the files
	if cfg.Simulate {
		newlineSW := new(bytes.Buffer)
		newlineSW.WriteString("# transformed files\n")
		for _, file := range setFiles {
			newlineSW.WriteString(fmt.Sprintf("## %s\n", file.Name))
			newlineSW.WriteString(fmt.Sprintf("```%s\n", file.Type))
			for _, line := range file.Lines {
				newlineSW.WriteString(line + "\n")
			}
			newlineSW.WriteString("\n```\n")
		}

		result := markdown.Render(newlineSW.String(), 132, 6)
		fmt.Println(string(result))
	} else {
		for _, file := range setFiles {
			logger.Info("writing changes to file %s", file.Path)

			outFile, err := os.OpenFile(file.Path, os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			defer outFile.Close()
			for _, line := range file.Lines {
				outFile.WriteString(line + "\n")
			}
		}
	}

	return nil
}

func setWellKnownValues(dictionary map[string]string) map[string]string {

	// set latest tag
	latest, err := getGitLatestTag()
	if err != nil {
		latest = "unknown"
	}
	dictionary["latesttag"] = latest

	// set actual date
	dictionary["actualdate"] = time.Now().Format("2006-01-02")

	// set acttual time
	dictionary["actualtime"] = time.Now().Format("15:04:05")

	//set actual timestamp
	dictionary["actualtimestamp"] = time.Now().Format("2006-01-02 15:04:05")

	return dictionary

}

func getGitLatestTag() (string, error) {
	gitRepo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	tags, err := gitRepo.Tags()
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`^\d\.`)

	versions := make([]string, 0)
	err = tags.ForEach(func(c *plumbing.Reference) error {
		v := c.Name().Short()
		if r.MatchString(v) {
			versions = append(versions, v)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Sort(byVersion(versions))
	if len(versions) > 0 {
		return versions[len(versions)-1], nil
	} else {
		return "unknown", nil
	}
}

type byVersion []string

func (s byVersion) Len() int {
	return len(s)
}

func (s byVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byVersion) Less(i, j int) bool {
	v1, err := version.NewVersion(s[i])
	if err != nil {
		panic(err)
	}
	v2, err := version.NewVersion(s[j])
	if err != nil {
		panic(err)
	}
	return v1.LessThan(v2)
}

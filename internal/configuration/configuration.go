package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// A placeholder is a comment line found in text in form of
// `$verifyr:template=<a golang template>$`
//
// The format of the comment may be different, since the source can be any textual format (.go, .yaml, .xml, .java, ...)
// The placeholder is used to identify the line where the template must be applied. The template is applied to the
// immediately following line.
//
// `Line` is the file line to be replaced by the result of the template application.
//
// `TemplateText` is the text of the template to be applied.
//
// `Template` is the parsed template.
type Placeholder struct {
	TemplateText string
	Template     *template.Template
	Line         int
}

func (p *Placeholder) String() string {
	res := "Placeholder{\n"
	res += "\tTemplateText:" + p.TemplateText + "\n"
	res += "\tLine:" + strconv.Itoa(p.Line) + "\n"
	res += "}"
	return res
}

// A ConfigFile is a file to be processed by versifyr.
//
// It contains the name of the file, the path to the file, the type of the file and the source of the file.
//
// Lines is the content of the file, as a slice of strings.
//
// The placeholders are the lines of the file where the templates must be applied. The placeholders define also the
// template to be applied.
type ConfigFile struct {
	Name         string         `koanf:"name"`
	Path         string         `koanf:"path"`
	Type         string         `koanf:"type"`
	Lines        []string       `koanf:"-"`
	Placeholders []*Placeholder `koanf:"-"`
}

func (c *ConfigFile) String() string {
	res := "ConfigFile{\n"
	res += "\tName: " + c.Name + "\n"
	res += "\tPath: " + c.Path + "\n"
	res += "\tType: " + c.Type + "\n"
	res += "\tPlaceholders:[\n"
	for i, p := range c.Placeholders {
		res += fmt.Sprintf("\t[%v]:%v\n", i, p)
		res += fmt.Sprintf("\ttarget => %v\n", c.Lines[p.Line])
	}
	res += "\t]\n"
	res += "}"
	return res
}

// Config is the configuration structure for the application.
type Config struct {
	Debug    bool          `koanf:"-"`
	BasePath string        `koanf:"-"`
	Simulate bool          `koanf:"-"`
	Files    []*ConfigFile `koanf:"files"`
}

func (c *Config) String() string {
	res := "Config{\n"
	for _, f := range c.Files {
		res += "\t" + f.String() + "\n"
	}
	res += "}"
	return res
}

// BASEPATH_DEFAULT is the default base path for the configuration file.
const BASEPATH_DEFAULT = ".versifyr"

// BASEPATH_ENV is the environment variable to override the default base path for the configuration file.
const BASEPATH_ENV = "VERSIFYR_BASEPATH"

// CONFIG_FILENAME is the name of the configuration file.
const CONFIG_FILENAME = "configuration.yaml"

// SAMPLE_CONFIG is the sample configuration file. It is used to create the configuration file with init comand.
const SAMPLE_CONFIG = `
# sample configuration file for versifyr
files:
  - name: version.go
    type: go
    path: internal/versifyr/version.go
	- name: chart.yaml
    type: yaml
    path: chart/Chart.yaml
  - name: pom.xml
    type: xml
    path: pom.xml
  - name: Version.java
    type: java
    path: src/main/java/sample/Version.java
`

// VERIFYER_TEMPLATE_START is the string that identifies the start of a template in a comment line.
const VERIFYER_TEMPLATE_START = "$versifyr:template="

// GetBasePath returns the base path for the configuration file. First default, then environment variable.
func GetBasePath() string {
	bp := BASEPATH_DEFAULT

	ebp := os.Getenv(BASEPATH_ENV)
	if ebp != "" {
		bp = ebp
	}

	return bp
}

// CreateConfiguration creates a sample configuration file.
func (cfg *Config) CreateConfiguration() error {
	if !cfg.Simulate {
		err := os.Mkdir(cfg.BasePath, 0755)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Simulating creation of configuration file - mkdir skipped")
		fmt.Println("folder: ", cfg.BasePath)
	}

	if !cfg.Simulate {
		fl, err := os.Create(cfg.BasePath + "/" + CONFIG_FILENAME)
		if err != nil {
			return err
		}
		defer fl.Close()

		_, err = fl.WriteString(SAMPLE_CONFIG)
		return err
	} else {
		fmt.Println("Simulating creation of configuration file - file creation skipped")
		fmt.Println("file: ", cfg.BasePath+"/"+CONFIG_FILENAME)
		fmt.Println("content:", SAMPLE_CONFIG)

		return nil
	}
}

// NewConfig creates a new configuration from the configuration file.
func NewConfig() (*Config, error) {
	var k = koanf.New(".")

	basepath := GetBasePath()
	c := &Config{
		Debug:    false,
		BasePath: basepath,
		Simulate: false,
	}

	err := k.Load(file.Provider(basepath+"/"+CONFIG_FILENAME), yaml.Parser())
	if err != nil {
		return c, err
	}

	err = k.Unmarshal("", &c)
	if err != nil {
		return nil, err
	}

	for _, f := range c.Files {
		f.Placeholders = make([]*Placeholder, 0)

		// load file
		bytes, err := os.ReadFile(f.Path)
		if err != nil {
			return c, err
		}

		content := string(bytes)

		// split lines
		lines := strings.Split(content, "\n")

		// find placeholders
		for i, l := range lines {
			if s := strings.Index(l, VERIFYER_TEMPLATE_START); s > -1 {
				e := strings.LastIndex(l, "$")
				if e == -1 || e < s+len(VERIFYER_TEMPLATE_START) {
					return nil, fmt.Errorf("template %s not closed at line %d", f.Name, i)
				}
				tpltxt := l[s+len(VERIFYER_TEMPLATE_START) : e]

				tpl, err := template.New(f.Name).Parse(tpltxt)
				if err != nil {
					return nil, err
				}

				ph := &Placeholder{
					TemplateText: tpltxt,
					Line:         i + 1,
					Template:     tpl,
				}
				f.Placeholders = append(f.Placeholders, ph)
			}
		}
		f.Lines = lines
	}
	return c, nil
}

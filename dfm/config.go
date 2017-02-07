package dfm

import (
	"bytes"
	"path"

	"github.com/benjamincaldwell/dfm/tasks"
	"github.com/benjamincaldwell/dfm/templates"
	"github.com/benjamincaldwell/go-printer"

	"github.com/benjamincaldwell/dfm/utilities"
	"github.com/ghodss/yaml"
)

type Configuration struct {
	Repo       string
	SrcDir     string
	DestDir    string
	Alias      map[string]string
	Tasks      map[string]tasks.Task
	homeDir    string
	configFile string
}

func (c *Configuration) Parse(file string) error {
	// data, err := ioutil.ReadFile(file)
	tpl, err := templates.New(
		templates.AppendFiles(file),
	)
	var data bytes.Buffer
	tpl.Execute(&data)
	if err == nil {
		err = yaml.Unmarshal(data.Bytes(), c)
		return err
	}
	return err
}

func parseConfig(file string) (*Configuration, error) {
	config := Configuration{}
	err := config.Parse(file)
	return &config, err
}

func (c *Configuration) SetDefaults() {
	if c.SrcDir == "" {
		c.SrcDir = path.Join(c.homeDir, ".dotfiles")
		printer.VerboseWarning("srcDir not specified. Defaulting to %s", c.SrcDir)
	}
	c.SrcDir = utilities.AbsPath(c.SrcDir, c.homeDir)

	if c.DestDir == "" {
		c.DestDir = c.homeDir
		printer.VerboseWarning("destDir not specified. Defaulting to %s", c.DestDir)
	}
	c.DestDir = utilities.AbsPath(c.DestDir, c.homeDir)
}

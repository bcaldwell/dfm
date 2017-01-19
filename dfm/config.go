package dfm

import (
	"io/ioutil"
	"path"

	"github.com/benjamincaldwell/dfm/tasks"
	"github.com/benjamincaldwell/go-printer"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Repo    string
	SrcDir  string
	DestDir string
	Alias   map[string]string
	Tasks   map[string]tasks.Task
}

func parseConfig(file string) (*Configuration, error) {
	config := Configuration{}
	data, err := ioutil.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(data, &config)
		return &config, err
	}
	return &config, err
}

func (c *Configuration) SetDefaults(homeDir string) {
	if c.SrcDir == "" {
		c.SrcDir = path.Join(homeDir, ".dotfiles")
		printer.VerboseWarning("srcDir not specified. Defaulting to %s", c.SrcDir)
	}

	if c.DestDir == "" {
		c.DestDir = homeDir
		printer.VerboseWarning("destDir not specified. Defaulting to %s", c.DestDir)
	}
}

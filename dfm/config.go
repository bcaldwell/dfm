package dfm

import (
	"io/ioutil"
	"path"

	"github.com/benjamincaldwell/dfm/tasks"
	"github.com/benjamincaldwell/go-printer"

	yaml "gopkg.in/yaml.v2"
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
	data, err := ioutil.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(data, c)
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

	if c.DestDir == "" {
		c.DestDir = c.homeDir
		printer.VerboseWarning("destDir not specified. Defaulting to %s", c.DestDir)
	}
}

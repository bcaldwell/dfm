package dfm

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Repo    string
	SrcDir  string
	DestDir string
	Alias   map[string]string
	Tasks   map[string]Task
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

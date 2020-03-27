package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ComposeConfig is the root of the Monako config
type ComposeConfig struct {
	BaseURL       string   `yaml:"baseURL"`
	Title         string   `yaml:"title"`
	Origins       []origin `yaml:"origins"`
	FileWhitelist []string `yaml:"whitelist"`
}

type origin struct {
	Source      string `yaml:"src"`
	Branch      string `yaml:"branch,omitempty"`
	EnvUsername string `yaml:"envusername,omitempty"`
	EnvPassword string `yaml:"envpassword,omitempty"`
	DirWithDocs string `yaml:"docdir,omitempty"`
	TargetDir   string `yaml:"targetdir,omitempty"`
}

// LoadConfig returns the Monako config from the given configfilepath
func LoadConfig(configfilepath string) (config ComposeConfig, err error) {

	source, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &config)
	return

}

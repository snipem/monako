package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var FileWhitelist = []string{".md", ".adoc", ".jpg", ".jpeg", ".svg", ".gif", ".png"}

type composeConfig struct {
	Source      string `yaml:"src"`
	Branch      string `yaml:"branch,omitempty"`
	EnvUsername string `yaml:"envusername,omitempty"`
	EnvPassword string `yaml:"envpassword,omitempty"`
	DirWithDocs string `yaml:"docdir,omitempty"`
	TargetDir   string `yaml:"targetdir,omitempty"`
}

func LoadConfig(configfilepath string) (config []composeConfig, err error) {

	source, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return nil, err
	}

	var out []composeConfig

	err = yaml.Unmarshal(source, &out)
	if err != nil {
		return nil, err
	}
	return out, nil

}

package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ComposeConfig struct {
	Source      string `yaml:"src"`
	Branch      string `yaml:"branch,omitempty"`
	EnvUsername string `yaml:"envusername,omitempty"`
	EnvPassword string `yaml:"envpassword,omitempty"`
	DirWithDocs string `yaml:"docdir,omitempty"`
	TargetDir   string `yaml:"targetdir,omitempty"`
}

func LoadConfig(configfilepath string) (config []ComposeConfig, err error) {

	source, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return nil, err
	}

	var out []ComposeConfig

	err = yaml.Unmarshal(source, &out)
	if err != nil {
		return nil, err
	}
	return out, nil

}

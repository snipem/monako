package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/snipem/monako/pkg/helpers"
	"gopkg.in/yaml.v2"
)

// ComposeConfig is the root of the Monako config
type ComposeConfig struct {
	BaseURL       string   `yaml:"baseURL"`
	Title         string   `yaml:"title"`
	Origins       []Origin `yaml:"origins"`
	FileWhitelist []string `yaml:"whitelist"`

	CompositionDir string
}

// Origin contains all information for a document origin
type Origin struct {
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

	// Set standard composition subdirectory
	config.CompositionDir = "compose/content/"
	return

}

// Compose builds the Monako directory structure
func (c *ComposeConfig) Compose() {

	for _, o := range c.Origins {
		g, fs := helpers.CloneDir(o.Source, o.Branch, os.Getenv(o.EnvUsername), os.Getenv(o.EnvPassword))
		helpers.CopyDir(g, fs, o.DirWithDocs, c.CompositionDir, c.FileWhitelist)
	}

}

// CleanUp removes the compose folder
func (c *ComposeConfig) CleanUp() {
	err := os.RemoveAll(c.CompositionDir)
	if err != nil {
		log.Fatalf("Error while cleaning up: %s", err)
	}
}

// SetTargetDir sets the target dir. Standard is relative to the current directory (".")
func (c *ComposeConfig) SetTargetDir(targetdir string) {
	c.CompositionDir = filepath.Join(targetdir, c.CompositionDir)
}

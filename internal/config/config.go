package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/snipem/monako/pkg/helpers"
	"gopkg.in/yaml.v2"
)

// ComposeConfig is the root of the Monako config
type ComposeConfig struct {
	BaseURL       string           `yaml:"baseURL"`
	Title         string           `yaml:"title"`
	Origins       []helpers.Origin `yaml:"origins"`
	FileWhitelist []string         `yaml:"whitelist"`

	CompositionDir string
	ContentDir     string
}

// LoadConfig returns the Monako config from the given configfilepath
func LoadConfig(configfilepath string) (config ComposeConfig, err error) {

	source, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &config)

	// Set standard composition subdirectory
	config.CompositionDir = "compose"
	config.ContentDir = "content"
	return

}

// Compose builds the Monako directory structure
func (c *ComposeConfig) Compose() {

	contentDir := filepath.Join(c.CompositionDir, c.ContentDir)

	for _, o := range c.Origins {
		if o.FileWhitelist == nil {
			o.FileWhitelist = c.FileWhitelist
		}
		o.CloneDir()
		o.ComposeDir(contentDir)
	}

}

// // CleanUp removes the compose folder
// func (c *ComposeConfig) CleanUp() {
// 	// err := os.RemoveAll(c.CompositionDir)
// 	// if err != nil {
// 	// 	log.Fatalf("Error while cleaning up: %s", err)
// 	// }
// }

// SetTargetDir sets the target dir. Standard is relative to the current directory (".")
func (c *ComposeConfig) SetTargetDir(targetdir string) {
	if targetdir != "" {
		c.CompositionDir = filepath.Clean(targetdir)
	}
}

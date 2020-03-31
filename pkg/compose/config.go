package compose

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config is the root of the Monako config
type Config struct {
	BaseURL       string   `yaml:"baseURL"`
	Title         string   `yaml:"title"`
	Origins       []Origin `yaml:"origins"`
	FileWhitelist []string `yaml:"whitelist"`

	// HugoWorkingDir is the working dir for the Composition
	HugoWorkingDir string

	// ContentWorkingDir is the main working dir and where all the content is stored in
	ContentWorkingDir string
}

// LoadConfig returns the Monako config from the given configfilepath
func LoadConfig(configfilepath string, targetdir string) (config Config, err error) {

	source, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &config)

	// Set standard composition subdirectory
	config.HugoWorkingDir = filepath.Join(targetdir, "compose")

	// As demanded by Hugo
	config.ContentWorkingDir = filepath.Join(config.HugoWorkingDir, "content")

	log.Fatal(config)
	return

}

// Compose builds the Monako directory structure
func (c *Config) Compose() {

	contentDir := filepath.Join(c.HugoWorkingDir, c.ContentWorkingDir)

	for _, o := range c.Origins {
		if o.FileWhitelist == nil {
			o.FileWhitelist = c.FileWhitelist
		}
		o.CloneDir()
		o.ComposeDir(contentDir)
	}

}

// CleanUp removes the compose folder
func (c *Config) CleanUp() {
	err := os.RemoveAll(c.HugoWorkingDir)
	if err != nil {
		log.Fatalf("Error while cleaning up: %s", err)
	}
}

// SetTargetDir sets the target dir. Standard is relative to the current directory (".")
func (c *Config) SetTargetDir(targetdir string) {
	if targetdir != "" {
		c.HugoWorkingDir = filepath.Clean(targetdir)
	}
}

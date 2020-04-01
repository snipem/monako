package compose

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// Config is the root of the Monako config
type Config struct {
	BaseURL       string   `yaml:"baseURL"`
	Title         string   `yaml:"title"`
	Origins       []Origin `yaml:"origins"`
	Logo          string   `yaml:"logo"`
	FileWhitelist []string `yaml:"whitelist"`

	// HugoWorkingDir is the working dir for the Composition
	HugoWorkingDir string

	// ContentWorkingDir is the main working dir and where all the content is stored in
	ContentWorkingDir string
}

// LoadConfig returns the Monako config from the given configfilepath
func LoadConfig(configfilepath string, workingdir string) (config Config, err error) {

	source, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &config)

	// Set standard composition subdirectory
	config.setWorkingDir(workingdir)

	// As demanded by Hugo
	config.ContentWorkingDir = filepath.Join(config.HugoWorkingDir, "content")

	// Set pointer to config for each origin
	for i := range config.Origins {
		config.Origins[i].config = &config
	}

	return

}

// Compose builds the Monako directory structure
func (config *Config) Compose() {

	// If Origin has now own whitelist, use the Compose Whitelist
	for _, o := range config.Origins {
		if o.FileWhitelist == nil {
			o.FileWhitelist = config.FileWhitelist
		}
		o.CloneDir()
		o.ComposeDir()
	}

}

// CleanUp removes the compose folder
func (config *Config) CleanUp() {

	if (config.HugoWorkingDir) == "." {
		log.WithFields(log.Fields{
			"Working dir": config.HugoWorkingDir,
			"BaseURL":     config.BaseURL,
		}).Fatalf("Hugo working dir can't be .")
	}
	err := os.RemoveAll(config.HugoWorkingDir)
	if err != nil {
		log.Fatalf("CleanUp: Error while cleaning up: %s", err)
	}
}

// setWorkingDir sets the target dir. Standard is relative to the current directory (".")
func (config *Config) setWorkingDir(targetdir string) {
	if targetdir != "" {
		config.HugoWorkingDir = filepath.Join(targetdir, "compose")
	}
}

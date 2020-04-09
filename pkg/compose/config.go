package compose

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/snipem/monako/internal/workarounds"
	"github.com/snipem/monako/pkg/helpers"
	"gopkg.in/yaml.v2"
)

// Config is the root of the Monako config
type Config struct {
	BaseURL       string   `yaml:"baseURL"`
	Title         string   `yaml:"title"`
	Origins       []Origin `yaml:"origins"`
	Logo          string   `yaml:"logo"`
	FileWhitelist []string `yaml:"whitelist"`

	DisableCommitInfo bool `yaml:"disableCommitInfo"`

	// HugoWorkingDir is the working dir for the Composition. For example "your/dir/compose"
	HugoWorkingDir string

	// ContentWorkingDir is the main working dir and where all the content is stored in. For example "your/dir/"
	ContentWorkingDir string
}

// CommandLineSettings contains all the flags and settings made via the command line in main
type CommandLineSettings struct {
	// ConfigFilePath is the path to the Monako config
	ConfigFilePath string
	// MenuConfigFilePath is the path to the Menu config
	MenuConfigFilePath string
	// ContentWorkingDir is the directory where files should be created. Home of the compose folder.
	ContentWorkingDir string
	// BaseURL is the BaseURL of the site
	BaseURL string
	// Trace activates function name based logging
	Trace bool
	// FailOnHugoError will fail Monako when there are Hugo errors during build
	FailOnHugoError bool
}

// LoadConfig returns the Monako config from the given configfilepath
func LoadConfig(configfilepath string, workingdir string) (config *Config, err error) {

	source, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return nil, err
	}

	config.initConfig(workingdir)

	return config, nil

}

// initConfig does necessary init steps on a newly created or read config
func (config *Config) initConfig(workingdir string) {

	// Set standard composition subdirectory
	config.setWorkingDir(workingdir)

	// As demanded by Hugo
	config.ContentWorkingDir = filepath.Join(config.HugoWorkingDir, "content")

	// Set pointer to config for each origin
	for i := range config.Origins {
		config.Origins[i].config = config
	}

}

// Compose builds the Monako directory structure
func (config *Config) Compose() {

	// If Origin has now own whitelist, use the Compose Whitelist
	for i := range config.Origins {
		if config.Origins[i].FileWhitelist == nil {
			config.Origins[i].FileWhitelist = config.FileWhitelist
		}
		config.Origins[i].CloneDir()
		config.Origins[i].ComposeDir()
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

	log.Infof("Cleaned up: %s", config.HugoWorkingDir)
}

// setWorkingDir sets the target dir. Standard is relative to the current directory (".")
func (config *Config) setWorkingDir(workingdir string) {
	if workingdir != "" {
		config.HugoWorkingDir = filepath.Join(workingdir, "compose")
	}
}

// Init loads the Monako config, adds Workarounds, cleans up the working dir,
// runs Hugo for initializing the workign dir
func Init(cliSettings CommandLineSettings) (config *Config) {

	config, err := LoadConfig(cliSettings.ConfigFilePath, cliSettings.ContentWorkingDir)
	if err != nil {
		log.Fatal(err)
	}

	// TODO Move to loadconfig parameters
	if cliSettings.BaseURL != "" {
		// Overwrite config base url if it is set as parameter
		config.BaseURL = cliSettings.BaseURL
	}

	if _, isGithub := os.LookupEnv("GITHUB_WORKFLOW"); isGithub {
		log.Warn("Don't apply workaround on Github Actions, cause its flaky")
	} else {
		workarounds.AddFakeAsciidoctorBinForDiagramsToPath(config.BaseURL)
	}

	config.CleanUp()

	err = helpers.HugoRun([]string{"--quiet", "new", "site", config.HugoWorkingDir})
	if err != nil {
		log.Fatal(err)
	}

	createMonakoStructureInHugoFolder(config, cliSettings.MenuConfigFilePath)

	return config

}

// Run runs Hugo with the composed Monako source
func (config *Config) Run() error {

	err := helpers.HugoRun([]string{"--source", config.HugoWorkingDir})
	if err != nil {
		return err
	}
	return nil
}

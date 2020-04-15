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
	FileBlacklist []string `yaml:"blacklist"`

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
	// ShowVersion shows the current version and exists
	ShowVersion bool
	// FailOnHugoError will fail Monako when there are Hugo errors during build
	FailOnHugoError bool
	// OnlyCompose will only compose files but not generate HTML
	OnlyCompose bool
	// OnlyGenerate will only generate HTML files but not compose them
	OnlyGenerate bool
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
		if config.Origins[i].FileBlacklist == nil {
			config.Origins[i].FileBlacklist = config.FileBlacklist
		}

		filesystem := config.Origins[i].CloneDir()
		config.Origins[i].ComposeDir(filesystem)

		// After processing the origin, delete repo for freeing up memory
		// containing the whole virtual filesystem. Can easily add up to
		// multiple gigabyte
		config.Origins[i].repo = nil

		// Performance analysis ------

		// Frees up some more megabyte
		// debug.FreeOSMemory()

		// if os.Getenv("MONAKO_LOG_HEAP") == "true" {

		// 	f, err := os.Create(filepath.Join(fmt.Sprintf("origin_%d.heap.fix.log", i)))
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	pprof.WriteHeapProfile(f)
		// 	f.Close()
		// }

		// End Performance analysis ------

	}

}

// CleanUp removes the compose folder
func (config *Config) CleanUp() {

	if (config.HugoWorkingDir) == "." {
		log.Fatalf("Hugo working dir can't be .")
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

// Init loads the Monako config, adds Workarounds, runs Hugo for initializing the working directory
func Init(cliSettings CommandLineSettings) (config *Config) {

	config, err := LoadConfig(cliSettings.ConfigFilePath, cliSettings.ContentWorkingDir)
	if err != nil {
		log.Fatal(err)
	}

	if cliSettings.BaseURL != "" {
		// Overwrite config base url if it is set as parameter
		config.BaseURL = cliSettings.BaseURL
	}

	if _, isGithub := os.LookupEnv("GITHUB_WORKFLOW"); isGithub {
		log.Warn("Don't apply workaround on Github Actions, because its flaky")
	} else {
		workarounds.AddFakeAsciidoctorBinForDiagramsToPath(config.BaseURL)
	}

	if !cliSettings.OnlyGenerate {
		// Dont do these steps if only generate
		config.CleanUp()

		createMonakoStructureInHugoFolder(config, cliSettings.MenuConfigFilePath)
	}

	return config

}

// Generate runs Hugo on the composed Monako source
func (config *Config) Generate() error {

	if _, err := os.Stat(config.HugoWorkingDir); os.IsNotExist(err) {
		log.Fatalf("%s does not exist, run monako only-generate before?", config.HugoWorkingDir)
	}

	err := helpers.HugoRun([]string{"--source", config.HugoWorkingDir})
	if err != nil {
		return err
	}
	return nil
}

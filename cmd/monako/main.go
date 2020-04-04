package main

// run: make run

import (
	"flag"

	"github.com/snipem/monako/pkg/compose"

	log "github.com/sirupsen/logrus"
)

// CommandLineFlags contains all the flags and settings made via the command line in main
type CommandLineFlags struct {
	ConfigFilePath     string
	MenuConfigFilePath string
	WorkingDir         string
	BaseURL            string
	Trace              bool
	FailOnHugoError    bool
}

func parseCommandLine() (cliSettings CommandLineFlags) {

	var configfilepath = flag.String("config", "config.monako.yaml", "Configuration file")
	var menuconfigfilepath = flag.String("menu-config", "config.menu.md", "Menu file for monako-book theme")
	var workingdir = flag.String("working-dir", ".", "Working dir for composed site")
	var baseURL = flag.String("base-url", "", "Custom base URL")
	var trace = flag.Bool("trace", false, "Enable trace logging")
	var failOnHugoError = flag.Bool("fail-on-error", false, "Fail on document conversion errors")

	flag.Parse()

	return CommandLineFlags{
		ConfigFilePath:     *configfilepath,
		MenuConfigFilePath: *menuconfigfilepath,
		WorkingDir:         *workingdir,
		BaseURL:            *baseURL,
		Trace:              *trace,
		FailOnHugoError:    *failOnHugoError,
	}
}

func main() {

	cliSettings := parseCommandLine()

	if cliSettings.Trace {
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config := compose.Init(
		cliSettings.ConfigFilePath,
		cliSettings.MenuConfigFilePath,
		cliSettings.WorkingDir,
		cliSettings.BaseURL)

	config.Compose()

	err := config.Run()
	if cliSettings.FailOnHugoError && err != nil {
		log.Fatal(err)
	}

}

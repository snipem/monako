package main

// run: make run

import (
	"flag"

	"github.com/snipem/monako/pkg/compose"

	log "github.com/sirupsen/logrus"
)

func parseCommandLine() (cliSettings compose.CommandLineSettings) {

	var configfilepath = flag.String("config", "config.monako.yaml", "Configuration file")
	var menuconfigfilepath = flag.String("menu-config", "config.menu.md", "Menu file for monako-book theme")
	var workingdir = flag.String("working-dir", ".", "Working dir for composed site")
	var baseURL = flag.String("base-url", "", "Custom base URL")
	var trace = flag.Bool("trace", false, "Enable trace logging")
	var failOnHugoError = flag.Bool("fail-on-error", false, "Fail on document conversion errors")

	flag.Parse()

	return compose.CommandLineSettings{
		ConfigFilePath:     *configfilepath,
		MenuConfigFilePath: *menuconfigfilepath,
		ContentWorkingDir:  *workingdir,
		BaseURL:            *baseURL,
		Trace:              *trace,
		FailOnHugoError:    *failOnHugoError,
	}
}

func main() {

	cliSettings := parseCommandLine()

	if cliSettings.Trace {
		log.SetLevel(logrus.DebugLevel)
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config := compose.Init(cliSettings)

	config.Compose()

	err := config.Run()
	if cliSettings.FailOnHugoError && err != nil {
		log.Fatal(err)
	}

}

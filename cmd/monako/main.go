package main

// run: make run

import (
	"flag"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/snipem/monako/internal/theme"
	"github.com/snipem/monako/internal/workarounds"
	"github.com/snipem/monako/pkg/compose"
	"github.com/snipem/monako/pkg/helpers"
)

func addWorkarounds(c *compose.Config) {
	if runtime.GOOS == "windows" {
		log.Println("Can't apply asciidoc diagram workaround on windows")
	} else {
		workarounds.AddFakeAsciidoctorBinForDiagramsToPath(c.BaseURL)
	}
}

func main() {

	var configfilepath = flag.String("config", "config.monako.yaml", "Configuration file")
	var menuconfigfilepath = flag.String("menu-config", "config.menu.md", "Menu file for monako-book theme")
	var workingdir = flag.String("working-dir", ".", "Working dir for composed site")
	var baseURLflag = flag.String("base-url", "", "Custom base URL")
	var trace = flag.Bool("trace", false, "Enable trace logging")
	var failOnError = flag.Bool("fail-on-error", false, "Fail on document conversion errors")

	flag.Parse()

	if *trace {
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config, err := compose.LoadConfig(*configfilepath, *workingdir)
	if err != nil {
		log.Fatal(err)
	}

	if *baseURLflag != "" {
		// Overwrite config base url if it is set as parameter
		config.BaseURL = *baseURLflag
	}

	addWorkarounds(config)

	config.CleanUp()

	err = helpers.HugoRun([]string{"--quiet", "new", "site", config.HugoWorkingDir})
	if *failOnError && err != nil {
		log.Fatal(err)
	}

	theme.CreateHugoPage(config, *menuconfigfilepath)

	config.Compose()

	err = helpers.HugoRun([]string{"--source", config.HugoWorkingDir})
	if *failOnError && err != nil {
		log.Fatal(err)
	}

}

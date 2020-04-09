package main

// run: make run

import (
	"flag"
	"os"

	"github.com/snipem/monako/pkg/compose"
	"github.com/snipem/monako/pkg/helpers"

	log "github.com/sirupsen/logrus"
)

func parseCommandLine() (cliSettings compose.CommandLineSettings) {

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var configfilepath = f.String("config", "config.monako.yaml", "Configuration file")
	var menuconfigfilepath = f.String("menu-config", "config.menu.md", "Menu file for monako-book theme")
	var workingdir = f.String("working-dir", ".", "Working dir for composed site")
	var baseURL = f.String("base-url", "", "Custom base URL")
	var trace = f.Bool("trace", false, "Enable trace logging")
	var failOnHugoError = f.Bool("fail-on-error", false, "Fail on document conversion errors")
	var onlyCompose = f.Bool("only-compose", false, "Only compose the Monako structure")
	var onlyGenerate = f.Bool("only-generate", false, "Only generate HTML files from an existing Monako structure")

	f.Parse(os.Args[1:])
	if *onlyCompose && *onlyGenerate {
		log.Fatal("only-compose and only-generate can't be set both")
	}

	return compose.CommandLineSettings{
		ConfigFilePath:     *configfilepath,
		MenuConfigFilePath: *menuconfigfilepath,
		ContentWorkingDir:  *workingdir,
		BaseURL:            *baseURL,
		Trace:              *trace,
		FailOnHugoError:    *failOnHugoError,
		OnlyCompose:        *onlyCompose,
		OnlyGenerate:       *onlyGenerate,
	}
}

func main() {

	cliSettings := parseCommandLine()

	if cliSettings.Trace {
		helpers.Trace()
	}

	config := compose.Init(cliSettings)
	if !cliSettings.OnlyGenerate {
		config.Compose()
	}

	if !cliSettings.OnlyCompose {
		err := config.Generate()
		if cliSettings.FailOnHugoError && err != nil {
			log.Fatal(err)
		}
	}

}

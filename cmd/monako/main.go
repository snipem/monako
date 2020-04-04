package main

// run: make run

import (
	"flag"

	"github.com/snipem/monako/pkg/compose"

	log "github.com/sirupsen/logrus"
)

func main() {

	var configfilepath = flag.String("config", "config.monako.yaml", "Configuration file")
	var menuconfigfilepath = flag.String("menu-config", "config.menu.md", "Menu file for monako-book theme")
	var workingdir = flag.String("working-dir", ".", "Working dir for composed site")
	var baseURLflag = flag.String("base-url", "", "Custom base URL")
	var trace = flag.Bool("trace", false, "Enable trace logging")
	var failOnHugoError = flag.Bool("fail-on-error", false, "Fail on document conversion errors")

	flag.Parse()

	if *trace {
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config := compose.Init(*configfilepath, *menuconfigfilepath, *workingdir, *baseURLflag)

	config.Compose()

	err := config.Run()
	if *failOnHugoError && err != nil {
		log.Fatal(err)
	}

}

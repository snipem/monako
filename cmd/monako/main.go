package main

// run: make run

import (
	"flag"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/snipem/monako/internal/config"
	"github.com/snipem/monako/internal/theme"
	"github.com/snipem/monako/internal/workarounds"
	"github.com/snipem/monako/pkg/helpers"
)

func compose(url string, branch string, subdir string, target string, username string, password string, whitelist []string) {

	fs := helpers.CloneDir(url, branch, username, password)
	helpers.CopyDir(fs, subdir, "compose/content/"+target+"/", whitelist)
}

func addWorkarounds() {
	if runtime.GOOS == "windows" {
		log.Println("Can't apply asciidoc diagram workaround on windows")
	} else {
		workarounds.AddFakeAsciidoctorBinForDiagramsToPath()
	}
}

func main() {

	var configfilepath = flag.String("config", "config.monako.yaml", "Configuration file")
	var menuconfigfilepath = flag.String("menu-config", "config.menu.md", "Menu file for monako-book theme")
	var trace = flag.Bool("trace", false, "Enable trace logging")

	flag.Parse()

	if *trace == true {
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config, err := config.LoadConfig(*configfilepath)
	if err != nil {
		log.Fatal(err)
	}

	helpers.CleanUp()
	addWorkarounds()

	helpers.HugoRun([]string{"--quiet", "new", "site", "compose"})
	theme.GetTheme(config, *menuconfigfilepath)

	for _, c := range config.Origins {
		compose(c.Source, c.Branch, c.DirWithDocs, c.TargetDir, os.Getenv(c.EnvUsername), os.Getenv(c.EnvPassword), config.FileWhitelist)
	}

	helpers.HugoRun([]string{"--source", "compose"})

}

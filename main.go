package main

// run: make run

import (
	"flag"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var fileWhitelist = []string{".md", ".adoc", ".jpg", ".jpeg", ".svg", ".gif", ".png"}

// Default file mode for temporary files

func compose(url string, branch string, subdir string, target string, username string, password string) {

	fs := cloneDir(url, branch, username, password)
	copyDir(fs, subdir, "compose/content/"+target+"/")
}

func addWorkarounds() {
	if runtime.GOOS == "windows" {
		log.Println("Can't apply asciidoc diagram workaround on windows")
	} else {
		addFakeAsciidoctorBinForDiagramsToPath()
	}
}

func main() {

	var configfilepath = flag.String("config", "config.monako.yaml", "Configuration file")
	var hugoconfigfilepath = flag.String("hugo-config", "config.hugo.toml", "Configuration file for Hugo")
	var menuconfigfilepath = flag.String("menu-config", "config.menu.md", "Menu file for monako-book theme")
	var trace = flag.Bool("trace", false, "Enable trace logging")

	flag.Parse()

	if *trace == true {
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config, err := loadConfig(*configfilepath)
	if err != nil {
		log.Fatal(err)
	}

	cleanUp()
	addWorkarounds()

	hugoRun([]string{"--quiet", "new", "site", "compose"})
	getTheme(*hugoconfigfilepath, *menuconfigfilepath)

	for _, c := range config {
		compose(c.Source, c.Branch, c.DirWithDocs, c.TargetDir, os.Getenv(c.EnvUsername), os.Getenv(c.EnvPassword))
	}

	hugoRun([]string{"--source", "compose"})

}

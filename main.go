package main

// run: make run

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var fileWhitelist = []string{".md", ".adoc", ".jpg", ".jpeg", ".svg", ".gif", ".png"}

// Default file mode for temporary files

func compose(url string, branch string, subdir string, target string, username string, password string) {

	fs := cloneDir(url, branch, username, password)
	copyDir(fs, subdir, "compose/content/docs/"+target+"/")
}

func main() {

	var configfilepath = flag.String("config", "config.yaml", "Configuration file, default: config.yaml")
	var hugoconfigfilepath = flag.String("hugo-config", "config.toml", "Configuration file for hugo, default: config.toml")
	var menuconfigfilepath = flag.String("menu-config", "index.md", "Menu file for hugo-book theme, default: index.md")
	var trace = flag.Bool("trace", false, "Enable trace logging")

	flag.Parse()

	if *trace == true {
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config, err := LoadConfig(*configfilepath)
	if err != nil {
		log.Fatal(err)
	}

	CleanUp()
	HugoRun([]string{"--quiet", "new", "site", "compose"})
	getTheme(*hugoconfigfilepath, *menuconfigfilepath)

	for _, c := range config {
		compose(c.Source, c.Branch, c.DirWithDocs, c.TargetDir, os.Getenv(c.EnvUsername), os.Getenv(c.EnvPassword))
	}

	HugoRun([]string{"--source", "compose"})

}

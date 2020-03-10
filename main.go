package main

// run: make run

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/codeskyblue/go-sh"
	"github.com/gohugoio/hugo/commands"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var fileWhitelist = []string{".md", ".adoc", ".jpg", ".jpeg", ".svg", ".gif"}

func cloneDir(url string, branch string, username string, password string) billy.Filesystem {

	log.Printf("Cloning in to %s with branch %s", url, branch)

	fs := memfs.New()

	basicauth := http.BasicAuth{}

	if username != "" && password != "" {
		log.Printf("Using username and password")
		basicauth = http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url,
		Depth:         1,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Auth:          &basicauth,
	})

	if err != nil {
		log.Fatal(err)
	}

	return fs
}

func cleanUp() {
	os.RemoveAll("compose")
}

func shouldIgnoreFile(filename string) bool {
	for _, whitelisted := range fileWhitelist {
		if strings.HasSuffix(strings.ToLower(filename), strings.ToLower(whitelisted)) {
			return false
		}
	}
	return true
}

func copyDir(fs billy.Filesystem, subdir string, target string) {

	log.Printf("Entering subdir %s of virtual filesystem from to target %s", subdir, target)
	var files []os.FileInfo

	fs, err := fs.Chroot(subdir)
	if err != nil {
		log.Fatal(err)
	}

	files, err = fs.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			// TODO is this memory consuming or is fsSubdir freed after recursion?
			// fsSubdir := fs
			copyDir(fs, file.Name(), target+file.Name())
			continue
		} else if shouldIgnoreFile(file.Name()) {
			continue
		}

		f, err := fs.Open(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		// TODO check if 0755 is good?
		err = os.MkdirAll(target, os.FileMode(0755))
		if err != nil {
			log.Fatal(err)
		}

		t, err := os.Create(target + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		if _, err = io.Copy(t, f); err != nil {
			log.Fatal(err)
		}

		log.Printf("Copied %s\n", file.Name())

	}

}

func hugoRun(args []string) {
	// args := []string{"--contentDir", "compose"}
	commands.Execute(args)
}

func compose(url string, branch string, subdir string, target string, username string, password string) {

	fs := cloneDir(url, branch, username, password)
	copyDir(fs, subdir, "compose/content/docs/"+target+"/")
}

func getTheme(hugoconfig string, menuconfig string) {
	// FIXME make me native
	sh.Command("wget", "-qO-", "https://github.com/alex-shpak/hugo-book/archive/v6.zip").Command("bsdtar", "-xvf-", "-C", "compose/themes/").Run()
	// TODO has to be TOML
	sh.Command("cp", hugoconfig, "compose/config.toml").Run()

	sh.Command("mkdir", "-p", "compose/content/menu/").Run()
	sh.Command("cp", menuconfig, "compose/content/menu/index.md").Run()

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

	cleanUp()
	hugoRun([]string{"--quiet", "new", "site", "compose"})
	getTheme(*hugoconfigfilepath, *menuconfigfilepath)

	for _, c := range config {
		compose(c.Source, c.Branch, c.DirWithDocs, c.TargetDir, os.Getenv(c.EnvUsername), os.Getenv(c.EnvPassword))
	}

	hugoRun([]string{"--source", "compose"})

}

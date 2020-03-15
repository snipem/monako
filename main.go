package main

// run: make run

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/artdarek/go-unzip"
	log "github.com/sirupsen/logrus"

	"github.com/codeskyblue/go-sh"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var fileWhitelist = []string{".md", ".adoc", ".jpg", ".jpeg", ".svg", ".gif", ".png"}

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

		var targetFilename = target + "/" + file.Name()

		if strings.HasSuffix(file.Name(), ".md") {
			var dirty, _ = ioutil.ReadAll(f)
			clean := MarkdownPostprocessing(dirty)
			ioutil.WriteFile(targetFilename, clean, os.FileMode(0755))
		} else if strings.HasSuffix(file.Name(), ".adoc") {
			var dirty, _ = ioutil.ReadAll(f)
			clean := AsciidocPostprocessing(dirty)
			ioutil.WriteFile(targetFilename, clean, os.FileMode(0755))
		} else {

			t, err := os.Create(targetFilename)
			if err != nil {
				log.Fatal(err)
			}

			if _, err = io.Copy(t, f); err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("Copied %s\n", file.Name())

		// s, _ := ioutil.ReadAll(f)
		// fmt.Println(string(s))

		// if err != nil {
		// 	log.Fatal(err)
		// }

	}

}

func compose(url string, branch string, subdir string, target string, username string, password string) {

	fs := cloneDir(url, branch, username, password)
	copyDir(fs, subdir, "compose/content/docs/"+target+"/")
}

func extractTheme() {
	themezip, err := Asset("tmp/theme.zip")
	if err != nil {
		log.Fatalf("Error loading theme %s", err)
	}

	// TODO Don't use local filesystem, keep it in memory
	tmpFile, err := ioutil.TempFile(os.TempDir(), "monako-theme-")
	if err != nil {
		fmt.Println("Cannot create temporary file", err)
	}
	tmpFile.Write(themezip)
	tempfilename := tmpFile.Name()

	// err = ioutil.WriteFile(tempfilename, themezip, os.FileMode(0755))
	if err != nil {
		log.Fatalf("Error writing temp theme %s", err)
	}

	// TODO Don't use a library that depends on local files
	uz := unzip.New(tempfilename, "compose/themes")
	err = uz.Extract()
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll(tempfilename)
}

func getTheme(hugoconfig string, menuconfig string) {

	extractTheme()
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

	CleanUp()
	HugoRun([]string{"--quiet", "new", "site", "compose"})
	getTheme(*hugoconfigfilepath, *menuconfigfilepath)

	for _, c := range config {
		compose(c.Source, c.Branch, c.DirWithDocs, c.TargetDir, os.Getenv(c.EnvUsername), os.Getenv(c.EnvPassword))
	}

	HugoRun([]string{"--source", "compose"})

}

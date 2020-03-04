package main

// run: make run

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/codeskyblue/go-sh"
	"github.com/gohugoio/hugo/commands"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func cloneDir(url string, branch string) billy.Filesystem {

	log.Printf("Cloning in to %s with branch %s", url, branch)

	fs := memfs.New()

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url,
		Depth:         1,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
	})

	if err != nil {
		log.Fatal(err)
	}

	return fs
}

func cleanUp() {
	os.RemoveAll("compose")
}

func copyDir(fs billy.Filesystem, subdir string, target string) {
	// TODO Dir not working

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
		// s, _ := ioutil.ReadAll(f)
		// fmt.Println(string(s))

		// if err != nil {
		// 	log.Fatal(err)
		// }

	}

	// CheckIfError(err)

}

func hugoRun(args []string) {
	// args := []string{"--contentDir", "compose"}
	commands.Execute(args)
}

func compose(url string, branch string, subdir string, target string) {

	fs := cloneDir(url, branch)
	copyDir(fs, subdir, "compose/content/"+target+"/")
}

func getTheme() {
	// FIXME make me native
	sh.Command("wget", "-qO-", "https://github.com/spf13/hyde/archive/master.zip").Command("bsdtar", "-xvf-", "-C", "compose/themes/").Run()
	sh.Command("echo", "theme = 'hyde-master'").Command("tee", "-a", "compose/config.toml").Run()

	sh.Command("cat", "sidebar.example.toml").Command("tee", "-a", "compose/config.toml").Run()
}

func main() {

	trace := true

	if trace == true {
		// Add line and filename to log
		log.SetReportCaller(true)
	}

	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cleanUp()
	hugoRun([]string{"new", "site", "compose"})
	getTheme()

	for _, c := range config {
		compose(c.Source, c.Branch, c.DirWithDocs, c.TargetDir)
	}

	hugoRun([]string{"--source", "compose"})

}

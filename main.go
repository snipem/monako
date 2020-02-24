package main

// run: make run

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/codeskyblue/go-sh"
	"github.com/gohugoio/hugo/commands"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func cloneDir(url string, branch string) billy.Filesystem {

	fmt.Printf("Cloning in to %s with branch %s, with dir %s", url, branch)

	fs := memfs.New()

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url,
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

	var files []os.FileInfo
	files, err := fs.ReadDir(subdir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}
		f, _ := fs.Open(file.Name())

		// TODO check if 0755 is good?
		err := os.MkdirAll(target, os.FileMode(0755))
		if err != nil {
			log.Fatal(err)
		}

		t, _ := os.Create(target + "/" + file.Name())

		if _, err = io.Copy(t, f); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", file.Name())
		s, _ := ioutil.ReadAll(f)
		fmt.Println(string(s))

		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("File contents: %s", content)
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

	cleanUp()
	hugoRun([]string{"new", "site", "compose"})
	getTheme()

	compose("https://github.com/snipem/commute-tube", "master", ".", "commute")
	compose("https://github.com/snipem/psnprices", "master", ".", "psnprices")

	hugoRun([]string{"--source", "compose", "serve"})

}

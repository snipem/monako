package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/snipem/monako/internal/workarounds"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var filemode = os.FileMode(0700)

func CloneDir(url string, branch string, username string, password string) billy.Filesystem {

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

func shouldIgnoreFile(filename string, whitelist []string) bool {
	for _, whitelisted := range whitelist {
		if strings.HasSuffix(strings.ToLower(filename), strings.ToLower(whitelisted)) {
			return false
		}
	}
	return true
}

func CopyDir(fs billy.Filesystem, subdir string, target string, whitelist []string) {

	log.Printf("Copying subdir '%s' to target dir %s", subdir, target)
	var files []os.FileInfo

	_, err := fs.Stat(subdir)
	if err != nil {
		log.Fatalf("Error while reading subdir %s: %s", subdir, err)
	}

	fs, err = fs.Chroot(subdir)
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
			CopyDir(fs, file.Name(), target+"/"+file.Name(), whitelist)
			continue
		} else if shouldIgnoreFile(file.Name(), whitelist) {
			continue
		}

		f, err := fs.Open(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(target, filemode)
		if err != nil {
			log.Fatal(err)
		}

		var targetFilename = target + "/" + file.Name()

		if strings.HasSuffix(file.Name(), ".md") {
			var dirty, _ = ioutil.ReadAll(f)
			clean := workarounds.MarkdownPostprocessing(dirty)
			ioutil.WriteFile(targetFilename, clean, filemode)
		} else if strings.HasSuffix(file.Name(), ".adoc") {
			var dirty, _ = ioutil.ReadAll(f)
			clean := workarounds.AsciidocPostprocessing(dirty)
			ioutil.WriteFile(targetFilename, clean, filemode)
		} else {

			t, err := os.Create(targetFilename)
			if err != nil {
				log.Fatal(err)
			}

			if _, err = io.Copy(t, f); err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("%s -> %s\n", file.Name(), targetFilename)

	}

}

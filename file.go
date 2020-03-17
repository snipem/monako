package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var filemode = os.FileMode(0700)

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
			copyDir(fs, file.Name(), target+"/"+file.Name())
			continue
		} else if shouldIgnoreFile(file.Name()) {
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
			clean := MarkdownPostprocessing(dirty)
			ioutil.WriteFile(targetFilename, clean, filemode)
		} else if strings.HasSuffix(file.Name(), ".adoc") {
			var dirty, _ = ioutil.ReadAll(f)
			clean := AsciidocPostprocessing(dirty)
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

		log.Printf("Copied %s\n", file.Name())

	}

}

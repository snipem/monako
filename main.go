package main

// run: make run

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

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

func composeDir(fs billy.Filesystem, dir string, target string) {
	// TODO Dir not working

	var files []os.FileInfo
	files, err := fs.ReadDir(dir)

	for _, file := range files {

		if file.IsDir() {
			continue
		}
		f, _ := fs.Open(file.Name())
		t, _ := os.Create("compose/" + file.Name())

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

func main() {

	fs := cloneDir("https://github.com/snipem/commute-tube", "master")
	composeDir(fs, ".", "target")
}

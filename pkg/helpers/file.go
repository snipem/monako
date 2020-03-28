package helpers

// run: make test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/snipem/monako/internal/workarounds"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var filemode = os.FileMode(0700)

// Asciidoc is a const for identifying Asciidoc Documents
const Asciidoc = "ASCIIDOC"

// Markdown is a const for identifying Markdown Documents
const Markdown = "MARKDOWN"

// CloneDir clones a HTTPS or lokal Git repository with the given branch and optional username and password.
// A virtual filesystem is returned containing the cloned files.
func CloneDir(url string, branch string, username string, password string) (*git.Repository, billy.Filesystem) {

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

	// TODO Check if we can check out less depth. Like depth = 1
	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url,
		Depth:         0, // problem with depth = 1 is that git log from older commits, can't be accessed
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Auth:          &basicauth,
	})

	if err != nil {
		log.Fatal(err)
	}

	return repo, fs
}

func shouldIgnoreFile(filename string, whitelist []string) bool {
	for _, whitelisted := range whitelist {
		if strings.HasSuffix(strings.ToLower(filename), strings.ToLower(whitelisted)) {
			return false
		}
	}
	return true
}

func isMarkdown(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".md"))
}

func isAsciidoc(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".adoc")) ||
		strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".asciidoc")) ||
		strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".asc"))
}

// DetermineFormat determines the markup format of a file by it's filename.
// Results can be Markdown and Asciidoc
func DetermineFormat(filename string) string {
	if isMarkdown(filename) {
		return Markdown
	} else if isAsciidoc(filename) {
		return Asciidoc
	} else {
		return ""
	}
}

// CopyDir copies a subdir of a virtual filesystem to a target in the local relative filesystem.
// The copied files can be limited by a whitelist. The Git repository is used to obtain Git commit
// information
func CopyDir(g *git.Repository, fs billy.Filesystem, source string, target string, whitelist []string) {

	source = filepath.Clean(source)
	target = filepath.Clean(target)

	log.Printf("Copying subdir '%s' to target dir '%s' ...", source, target)

	var err error

	// TODO: This is also done on every recursion. Maybe this is overhead.
	for _, dir := range strings.Split(source, string(filepath.Separator)) {
		fs, err = fs.Chroot(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	var files []os.FileInfo
	files, err = fs.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			foldername := filepath.Join(target, file.Name())
			// TODO is this memory consuming or is fsSubdir freed after recursion?
			// fsSubdir := fs
			CopyDir(g, fs, filepath.Join(source, file.Name()), foldername, whitelist)
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

		var targetFilename = filepath.Join(target, file.Name())
		contentFormat := DetermineFormat(file.Name())

		gitFilepath, _ := filepath.Rel("/", filepath.Join(fs.Root(), file.Name()))

		switch contentFormat {
		case Asciidoc, Markdown:

			// TODO: Only use strings not []byte
			commitinfo, err := GetCommitInfo(g, gitFilepath)
			if err != nil {
				log.Fatal(err)
			}

			var dirty, _ = ioutil.ReadAll(f)
			var content []byte

			if contentFormat == Markdown {
				content = workarounds.MarkdownPostprocessing(dirty)
			} else if contentFormat == Asciidoc {
				content = workarounds.AsciidocPostprocessing(dirty)
			}
			content = []byte(ExpandFrontmatter(string(content), commitinfo))
			ioutil.WriteFile(targetFilename, content, filemode)

		default:
			copyFile(targetFilename, f)
		}

		log.Printf("%s -> %s\n", gitFilepath, targetFilename)

	}

}

func copyFile(targetFilename string, from io.Reader) {
	t, err := os.Create(targetFilename)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(t, from); err != nil {
		log.Fatal(err)
	}

}

// GetCommitInfo returns the Commit Info for a given file of the repository
// identified by it's filename
func GetCommitInfo(r *git.Repository, filename string) (*object.Commit, error) {

	cIter, err := r.Log(&git.LogOptions{
		FileName: &filename,
		All:      true,
	})

	if err != nil {
		return nil, fmt.Errorf("Error while opening %s from git log: %s", filename, err)
	}

	var commit *object.Commit

	commit, err = cIter.Next()
	defer cIter.Close()

	if err != nil {
		if err.Error() == "EOF" {
			return nil, fmt.Errorf("File %s not found in git log", filename)
		}
		return nil, fmt.Errorf("Unknown error while fetching git commit info for '%s' from git log", filename)
	}

	return commit, nil
}

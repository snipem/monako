package compose

// run: make test

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/snipem/monako/pkg/helpers"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Asciidoc is a const for identifying Asciidoc Documents
const Asciidoc = "ASCIIDOC"

// Markdown is a const for identifying Markdown Documents
const Markdown = "MARKDOWN"

// CloneDir clones a HTTPS or lokal Git repository with the given branch and optional username and password.
// A virtual filesystem is returned containing the cloned files.
func (origin *Origin) CloneDir() {

	log.Printf("\nCloning in to '%s' with branch '%s' ...\n", origin.URL, origin.Branch)

	origin.filesystem = memfs.New()

	basicauth := http.BasicAuth{}

	username := os.Getenv(origin.EnvUsername)
	password := os.Getenv(origin.EnvPassword)

	if username != "" && password != "" {
		log.Printf("Using username and password stored in env variables\n")
		basicauth = http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	// TODO Check if we can check out less depth. Like depth = 1
	repo, err := git.Clone(memory.NewStorage(), origin.filesystem, &git.CloneOptions{
		URL:           origin.URL,
		Depth:         0, // problem with depth = 1 is that git log from older commits, can't be accessed
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", origin.Branch)),
		SingleBranch:  true,
		Auth:          &basicauth,
	})

	if err != nil {
		log.Fatal(err)
	}

	origin.repo = repo

}

// GetFormat determines the markup format of a file by it's filename.
// Results can be Markdown and Asciidoc
func (file OriginFile) GetFormat() string {
	if helpers.IsMarkdown(file.RemotePath) {
		return Markdown
	} else if helpers.IsAsciidoc(file.RemotePath) {
		return Asciidoc
	} else {
		return ""
	}
}

// Origin contains all information for a document origin
type Origin struct {
	URL           string   `yaml:"src"`
	Branch        string   `yaml:"branch,omitempty"`
	EnvUsername   string   `yaml:"envusername,omitempty"`
	EnvPassword   string   `yaml:"envpassword,omitempty"`
	SourceDir     string   `yaml:"docdir,omitempty"`
	TargetDir     string   `yaml:"targetdir,omitempty"`
	FileWhitelist []string `yaml:"whitelist,omitempty"`

	Files []OriginFile

	repo       *git.Repository
	config     *Config
	filesystem billy.Filesystem
}

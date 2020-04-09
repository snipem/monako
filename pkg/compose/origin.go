package compose

// run: make test

import (
	"fmt"
	"os"
	"path"

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

	fmt.Printf("\nCloning in to '%s' with branch '%s' ...\n", origin.URL, origin.Branch)
	log.Debugf("Start cloning of %s", origin.URL)

	origin.filesystem = memfs.New()

	basicauth := http.BasicAuth{}

	username := os.Getenv(origin.EnvUsername)
	password := os.Getenv(origin.EnvPassword)

	if username != "" && password != "" {
		fmt.Printf("Using username and password stored in env variables\n")
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
	log.Debugf("Start cloning of %s", origin.URL)

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

// ComposeDir copies a subdir of a virtual filesystem to a target in the local relative filesystem.
// The copied files can be limited by a whitelist. The Git repository is used to obtain Git commit
// information
func (origin *Origin) ComposeDir() {
	origin.Files = origin.getWhitelistedFiles(origin.SourceDir)

	if len(origin.Files) == 0 {
		log.Printf("Found no matching files in '%s' with branch '%s' in folder '%s'\n", origin.URL, origin.Branch, origin.SourceDir)
	}

	for _, file := range origin.Files {
		file.composeFile()
	}
}

// NewOrigin returns a new origin with all needed fields
func NewOrigin(url string, branch string, sourceDir string, targetDir string) *Origin {
	o := new(Origin)
	o.URL = url
	o.Branch = branch
	o.SourceDir = sourceDir
	o.TargetDir = targetDir
	return o
}

func (origin *Origin) getWhitelistedFiles(startdir string) []OriginFile {

	var originFiles []OriginFile

	files, _ := origin.filesystem.ReadDir(startdir)
	for _, file := range files {

		// This is the path as stored in the remote repo
		// This can only be gathered here, because of recursing through
		// the file system
		// Use path here to support unixoid Git paths
		remotePath := path.Join(startdir, file.Name())

		if file.IsDir() {
			// Recurse over file and add their files to originFiles
			originFiles = append(
				originFiles,
				origin.getWhitelistedFiles(
					remotePath,
				)...)
		} else if helpers.FileIsWhitelisted(file.Name(), origin.FileWhitelist) {

			// Add the current file to the list of files returned
			originFiles = append(
				originFiles,
				origin.newFile(remotePath))
		}

	}
	return originFiles
}

func (origin *Origin) newFile(remotePath string) OriginFile {
	localPath := getLocalFilePath(origin.config.ContentWorkingDir, origin.SourceDir, origin.TargetDir, remotePath)

	originFile := OriginFile{
		RemotePath: remotePath,
		LocalPath:  localPath,

		parentOrigin: origin,
	}

	if !origin.config.DisableCommitInfo {

		commitinfo, err := originFile.getCommitInfo()
		if err != nil {
			log.Warnf("Can't extract Commit Info for '%s'", err)
		}
		originFile.Commit = commitinfo
	}

	return originFile
}

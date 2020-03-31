package compose

// run: make test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/snipem/monako/internal/workarounds"
	"github.com/snipem/monako/pkg/helpers"
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
func (origin *Origin) CloneDir() {

	log.Printf("Cloning in to %s with branch %s", origin.URL, origin.Branch)

	origin.filesystem = memfs.New()

	basicauth := http.BasicAuth{}

	username := os.Getenv(origin.EnvUsername)
	password := os.Getenv(origin.EnvPassword)

	if username != "" && password != "" {
		log.Printf("Using username and password")
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

	origin.repo = repo

	if err != nil {
		log.Fatal(err)
	}

	return
}

// GetFormat determines the markup format of a file by it's filename.
// Results can be Markdown and Asciidoc
func (file OriginFile) GetFormat() string {
	if helpers.IsMarkdown(file.Path) {
		return Markdown
	} else if helpers.IsAsciidoc(file.Path) {
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
	filesystem billy.Filesystem
}

type OriginFile struct {
	Commit *object.Commit
	Path   string

	parentOrigin *Origin
}

// ComposeDir copies a subdir of a virtual filesystem to a target in the local relative filesystem.
// The copied files can be limited by a whitelist. The Git repository is used to obtain Git commit
// information
func (origin Origin) ComposeDir(rootDir string) {
	origin.Files = origin.getWhitelistedFiles(origin.SourceDir)

	for _, file := range origin.Files {
		file.composeFile(rootDir)
	}
}

func NewOrigin(url string, branch string, sourceDir string, targetDir string) *Origin {
	o := new(Origin)
	o.URL = url
	o.Branch = branch
	o.SourceDir = sourceDir
	o.TargetDir = targetDir
	return o
}

func (origin Origin) getWhitelistedFiles(startdir string) []OriginFile {

	var originFiles []OriginFile

	files, _ := origin.filesystem.ReadDir(startdir)
	for _, file := range files {

		if file.IsDir() {
			// Recurse over file and add there files to originFiles
			originFiles = append(originFiles,
				origin.getWhitelistedFiles(
					filepath.Join(startdir, file.Name()),
				)...)
		} else if helpers.FileIsWhitelisted(file.Name(), origin.FileWhitelist) {
			// Just add this file to originFiles
			var originFile OriginFile
			originFile.Path = filepath.Join(startdir, file.Name())
			originFile.parentOrigin = &origin

			// var err error
			// originFile.Commit, err = GetCommitInfo(origin.repo, originFile.Path)

			// if err != nil {
			// 	log.Fatalf("Can't extract git info for %s: %s", originFile.Path, err)
			// }

			log.Println(originFile.Path)

			originFiles = append(originFiles, originFile)
		}

	}
	return originFiles
}

func (file OriginFile) composeFile(rootDir string) {

	sourceDir := file.parentOrigin.SourceDir
	relativeFilePath := strings.TrimPrefix(file.Path, sourceDir)

	fileDirs := filepath.Dir(relativeFilePath)
	copyDir := filepath.Join(rootDir, file.parentOrigin.TargetDir, fileDirs)

	// log.Printf("Trying to create '%s'", copyDir)
	err := os.MkdirAll(copyDir, filemode)
	if err != nil {
		log.Fatalf("Error when creating '%s': %s", copyDir, err)
	}

	var targetFilename = filepath.Join(rootDir, file.parentOrigin.TargetDir, relativeFilePath)
	contentFormat := file.GetFormat()

	// gitFilepath, _ := filepath.Rel("/", filepath.Join(fs.Root(), file.Name()))

	switch contentFormat {
	case Asciidoc, Markdown:
		file.copyMarkupFile(targetFilename)
	default:
		file.copyRegularFile(targetFilename)
	}
	log.Printf("%s -> %s\n", file.Path, targetFilename)

}

func (file OriginFile) copyRegularFile(targetFilename string) {

	origin := file.parentOrigin
	f, err := origin.filesystem.Open(file.Path)

	if err != nil {
		log.Fatal(err)
	}
	t, err := os.Create(targetFilename)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(t, f); err != nil {
		log.Fatal(err)
	}

}

func (file OriginFile) copyMarkupFile(targetFilename string) {

	// TODO: Only use strings not []byte
	// commitinfo, err := GetCommitInfo(g, gitFilepath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	bf, err := file.parentOrigin.filesystem.Open(file.Path)

	var dirty, _ = ioutil.ReadAll(bf)
	var content []byte
	contentFormat := file.GetFormat()

	if contentFormat == Markdown {
		content = workarounds.MarkdownPostprocessing(dirty)
	} else if contentFormat == Asciidoc {
		content = workarounds.AsciidocPostprocessing(dirty)
	}
	// content = []byte(ExpandFrontmatter(string(content), gitFilepath, file.Commit))
	err = ioutil.WriteFile(targetFilename, content, filemode)
	if err != nil {
		log.Fatalf("Error writing file %s", err)
	}
}

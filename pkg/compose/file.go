package compose

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
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// OriginFile represents a single file of an origin
type OriginFile struct {
	Commit *object.Commit
	// RemotePath is the path in the origin repository
	RemotePath string
	// LocalPath is the absolute path on the local disk
	LocalPath string

	parentOrigin *Origin
}

// ComposeDir copies a subdir of a virtual filesystem to a target in the local relative filesystem.
// The copied files can be limited by a whitelist. The Git repository is used to obtain Git commit
// information
func (origin Origin) ComposeDir() {
	origin.Files = origin.getWhitelistedFiles(origin.SourceDir)

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
			originFile.RemotePath = filepath.Join(startdir, file.Name())

			originFile.parentOrigin = &origin

			originFile.LocalPath = getLocalFilePath(origin.config.ContentWorkingDir, origin.SourceDir, origin.TargetDir, originFile.RemotePath)
			log.Info(originFile)

			// var err error
			// originFile.Commit, err = GetCommitInfo(origin.repo, originFile.Path)

			// if err != nil {
			// 	log.Fatalf("Can't extract git info for %s: %s", originFile.Path, err)
			// }

			log.Println(originFile.RemotePath)

			originFiles = append(originFiles, originFile)
		}

	}
	return originFiles
}

func (file OriginFile) composeFile() {

	file.createParentDir()
	contentFormat := file.GetFormat()

	switch contentFormat {
	case Asciidoc, Markdown:
		file.copyMarkupFile()
	default:
		file.copyRegularFile()
	}
	fmt.Printf("%s -> %s\n", file.RemotePath, file.LocalPath)

}

// createParentDir creates the parent directories for the file in the local filesystem
func (file OriginFile) createParentDir() {
	log.Debugf("Creating local folder '%s'", filepath.Dir(file.LocalPath))
	err := os.MkdirAll(filepath.Dir(file.LocalPath), filemode)
	if err != nil {
		log.Fatalf("Error when creating '%s': %s", filepath.Dir(file.LocalPath), err)
	}
}

func (file OriginFile) copyRegularFile() {

	origin := file.parentOrigin
	f, err := origin.filesystem.Open(file.RemotePath)

	if err != nil {
		log.Fatal(err)
	}
	t, err := os.Create(file.LocalPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(t, f); err != nil {
		log.Fatal(err)
	}

}

func (file OriginFile) copyMarkupFile() {

	// TODO: Only use strings not []byte
	// commitinfo, err := GetCommitInfo(g, gitFilepath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	bf, err := file.parentOrigin.filesystem.Open(file.RemotePath)

	var dirty, _ = ioutil.ReadAll(bf)
	var content []byte
	contentFormat := file.GetFormat()

	if contentFormat == Markdown {
		content = workarounds.MarkdownPostprocessing(dirty)
	} else if contentFormat == Asciidoc {
		content = workarounds.AsciidocPostprocessing(dirty)
	}
	// content = []byte(ExpandFrontmatter(string(content), gitFilepath, file.Commit))
	err = ioutil.WriteFile(file.LocalPath, content, filemode)
	if err != nil {
		log.Fatalf("Error writing file %s", err)
	}
}

func getLocalFilePath(composeDir, remoteDocDir string, targetDir string, remoteFile string) string {
	// Since a remoteDocDir is defined, this should not be created in the local filesystem
	relativeFilePath := strings.TrimPrefix(remoteFile, remoteDocDir)
	return filepath.Join(composeDir, targetDir, relativeFilePath)
}

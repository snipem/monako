package compose

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/parser/pageparser"
	log "github.com/sirupsen/logrus"

	"github.com/snipem/monako/internal/workarounds"
	"github.com/snipem/monako/pkg/helpers"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/yaml.v2"
)

const standardFilemode = os.FileMode(0700)

// OriginFile represents a single file of an origin
type OriginFile struct {

	// Commit is the commit info about this file
	Commit *object.Commit
	// RemotePath is the path in the origin repository
	RemotePath string
	// LocalPath is the absolute path on the local disk
	LocalPath string

	// parentOrigin of this file
	parentOrigin *Origin
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
		remotePath := filepath.Join(startdir, file.Name())

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

	commitinfo, err := originFile.getCommitInfo()
	if err != nil {
		log.Warnf("Can't extract Commit Info for '%s'", err)
	}

	originFile.Commit = commitinfo

	return originFile
}

func (file *OriginFile) composeFile() {

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
func (file *OriginFile) createParentDir() {
	log.Debugf("Creating local folder '%s'", filepath.Dir(file.LocalPath))
	err := os.MkdirAll(filepath.Dir(file.LocalPath), standardFilemode)
	if err != nil {
		log.Fatalf("Error when creating '%s': %s", filepath.Dir(file.LocalPath), err)
	}
}

func (file *OriginFile) copyRegularFile() {

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

// getCommitInfo returns the Commit Info for a given file of the repository
// identified by it's filename
func (file *OriginFile) getCommitInfo() (*object.Commit, error) {

	r := file.parentOrigin.repo
	cIter, err := r.Log(&git.LogOptions{
		FileName: &file.RemotePath,
		Order:    git.LogOrderCommitterTime,
	})

	if err != nil {
		return nil, fmt.Errorf("Error while opening %s from git log: %s", file.RemotePath, err)
	}

	returnCommit, err := cIter.Next()

	if err != nil {
		return nil, fmt.Errorf("File not found in git log: '%s'", file.RemotePath)
	}

	log.Debugf("Git Commit found for %s, %s", file.RemotePath, returnCommit)

	// This has to be here, otherwise the iterator will return garbage
	defer cIter.Close()

	return returnCommit, nil
}

func (file *OriginFile) copyMarkupFile() {

	// TODO: Only use strings not []byte

	bf, err := file.parentOrigin.filesystem.Open(file.RemotePath)
	if err != nil {
		log.Fatalf("Error copying markup file %s", err)
	}

	var dirty, _ = ioutil.ReadAll(bf)
	var content []byte
	contentFormat := file.GetFormat()

	if contentFormat == Markdown {
		content = workarounds.MarkdownPostprocessing(dirty)
	} else if contentFormat == Asciidoc {
		content = workarounds.AsciidocPostprocessing(dirty)
	}

	// TODO: Add ExpandFrontmatter function
	content = []byte(file.ExpandFrontmatter(string(content)))

	err = ioutil.WriteFile(file.LocalPath, content, standardFilemode)
	if err != nil {
		log.Fatalf("Error writing file %s", err)
	}
}

// getLocalFilePath returns the desired local file path for a remote file in the local filesystem.
// It is based on the local absolute composeDir, the remoteDocDir to strip it's path from the local file,
// the target dir to generate the local path and the file name itself
func getLocalFilePath(composeDir, remoteDocDir string, targetDir string, remoteFile string) string {
	// Since a remoteDocDir is defined, this should not be created in the local filesystem
	relativeFilePath := strings.TrimPrefix(remoteFile, remoteDocDir)
	return filepath.Join(composeDir, targetDir, relativeFilePath)
}

// ExpandFrontmatter expands the existing frontmatter with the parameters given
func (file *OriginFile) ExpandFrontmatter(content string) string {

	if file.Commit == nil {
		log.Debug("Git Info is not set, returning without adding it")
		return content
	}

	oldFrontmatter, body := splitFrontmatterAndBody(content)

	return fmt.Sprintf(`---
%s

MonakoGitRemote: %s
MonakoGitRemotePath: %s
MonakoGitURL: %s
MonakoGitLastCommitDate: %s
MonakoGitLastCommitAuthor: %s
MonakoGitLastCommitAuthorEmail: %s
---

`+body,
		oldFrontmatter,
		file.parentOrigin.URL,
		file.RemotePath,
		getWebLinkForGit(
			file.parentOrigin.URL,
			file.parentOrigin.Branch,
			file.RemotePath,
		),
		file.Commit.Author.When.Format(time.RFC3339),
		file.Commit.Author.Name,
		file.Commit.Author.Email)

}

func splitFrontmatterAndBody(content string) (frontmatter string, body string) {
	// TODO Convert from toml, yaml, etc
	contentFrontmatter, err := pageparser.ParseFrontMatterAndContent(strings.NewReader(content))
	if err != nil {
		log.Fatalf("Error while splitting frontmatter: %s", err)
	}
	contentMarshaled, err := yaml.Marshal(contentFrontmatter.FrontMatter)

	return string(contentMarshaled), string(contentFrontmatter.Content)
}

func getWebLinkForGit(gitURL string, branch string, remotePath string) string {

	// TODO Maybe return nothing if it's a ssh or file repository
	// URLs for checkout have .git suffix
	gitURL = strings.TrimSuffix(gitURL, ".git")
	u, err := url.Parse(gitURL)
	if err != nil {
		log.Fatalf("Can't parse url: %s", gitURL)
	}
	u.Path = path.Join(u.Path, "blob", branch, remotePath)
	return u.String()
}

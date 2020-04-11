package compose

// run: MONAKO_HUGE_REPOS_TEST=true go test -v ./pkg/compose/ -run TestHugeRepositories

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
func getCommitInfo(remotePath string, repo *git.Repository) (*object.Commit, error) {

	log.Debugf("Getting commit info for %s", remotePath)

	if repo == nil {
		return nil, fmt.Errorf("Repository is nil")
	}

	// Problem seems to be the longer the file hasn't been in the log, the longer it takes to retrieve it
	cIter, err := repo.Log(&git.LogOptions{
		FileName: &remotePath,
		Order:    git.LogOrderCommitterTime,
	})

	if err != nil {
		return nil, fmt.Errorf("Error while opening %s from git log: %s", remotePath, err)
	}

	returnCommit, err := cIter.Next()

	if err != nil {
		return nil, fmt.Errorf("File not found in git log: '%s'", remotePath)
	}

	log.Debugf("Git Commit found for %s, %s", remotePath, returnCommit)

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
MonakoGitLastCommitHash: %s
MonakoGitURLCommit: %s
lastMod: %s
MonakoGitLastCommitAuthor: %s
MonakoGitLastCommitAuthorEmail: %s
---

`+body,
		oldFrontmatter,
		file.parentOrigin.URL,
		file.RemotePath,
		getWebLinkForFileInGit(
			file.parentOrigin.URL,
			file.parentOrigin.Branch,
			file.RemotePath,
		),
		file.Commit.Hash,
		getWebLinkForGitCommit(
			file.parentOrigin.URL,
			file.Commit.Hash.String(),
		),
		// Use lastMod because other variables won't be parsed as date by Hugo
		// Resulting in no date format functions on the file
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

	// No frontmatter found, return old content
	if contentFrontmatter.Content == nil {
		return "", content
	}

	contentMarshaled, err := yaml.Marshal(contentFrontmatter.FrontMatter)
	if err != nil {
		log.Fatalf("Error while marshalling frontmatter to YAML: %s", err)
	}

	return string(contentMarshaled), string(contentFrontmatter.Content)
}

func getWebLinkForFileInGit(gitURL string, branch string, remotePath string) string {

	// URLs for checkout have .git suffix
	gitURL = strings.TrimSuffix(gitURL, ".git")
	u, err := url.Parse(gitURL)

	if err != nil {
		return ""
	}

	if !strings.Contains(u.Scheme, "http") {
		return ""
	}

	// This works for Github and Gitlab
	middlePath := "blob"

	// Bitbucket does it differently
	if strings.Contains(u.Host, "bitbucket") {
		middlePath = "src"
	}

	u.Path = path.Join(u.Path, middlePath, branch, remotePath)
	return u.String()
}

func getWebLinkForGitCommit(gitURL string, commitID string) string {
	// URLs for checkout have .git suffix
	gitURL = strings.TrimSuffix(gitURL, ".git")
	u, err := url.Parse(gitURL)

	if err != nil {
		return ""
	}

	if !strings.Contains(u.Scheme, "http") {
		return ""
	}

	// This works for Github and Gitlab
	middlePath := "commit"

	// Bitbucket does it differently
	if strings.Contains(u.Host, "bitbucket") {
		middlePath = "commits"
	}

	u.Path = path.Join(u.Path, middlePath, commitID)
	return u.String()
}

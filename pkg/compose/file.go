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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/snipem/monako/internal/workarounds"
	"github.com/snipem/monako/pkg/helpers"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

const standardFilemode = os.FileMode(0700)

// OriginFile represents a single file of an origin
type OriginFile struct {

	// Commit is the commit info about this file
	Commit *OriginFileCommit
	// RemotePath is the path in the origin repository
	RemotePath string
	// LocalPath is the absolute path on the local disk
	LocalPath string

	// parentOrigin of this file
	parentOrigin *Origin
}

// OriginFileCommit represents a commit
type OriginFileCommit struct {
	Hash   string
	Author OriginFileCommitter
	Date   time.Time
}

// OriginFileCommitter represents the committer of a commit
type OriginFileCommitter struct {
	Name  string
	Email string
}

func (file *OriginFile) composeFile(filesystem billy.Filesystem) error {

	err := createParentDir(file.LocalPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error composing file %s", file.LocalPath))
	}
	contentFormat := file.GetFormat()

	switch contentFormat {
	case Asciidoc, Markdown:
		err := file.copyMarkupFile(filesystem)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error copying markup file"))
		}
	default:
		err := file.copyRegularFile(filesystem)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error copying regular file"))
		}
	}
	fmt.Printf("%s -> %s\n", file.RemotePath, file.LocalPath)
	return nil

}

// GetFormat determines the markup format of a file by it's filename.
// Results can be Markdown and Asciidoc
func (file *OriginFile) GetFormat() string {
	if helpers.IsMarkdown(file.RemotePath) {
		return Markdown
	} else if helpers.IsAsciidoc(file.RemotePath) {
		return Asciidoc
	} else {
		return ""
	}
}

// createParentDir creates the parent directories for the file in the local filesystem
func createParentDir(localPath string) error {
	log.Debugf("Creating parent dir '%s'", filepath.Dir(localPath))
	err := os.MkdirAll(filepath.Dir(localPath), standardFilemode)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating parent dir %s", localPath))
	}
	return nil
}

func (file *OriginFile) copyRegularFile(filesystem billy.Filesystem) error {

	f, err := filesystem.Open(file.RemotePath)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error opening regular remote file for copying %s", file.RemotePath))
	}
	t, err := os.Create(file.LocalPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating regular local file %s", file.LocalPath))
	}

	if _, err = io.Copy(t, f); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error copying regular remote file to local file %s -> %s", file.RemotePath, file.LocalPath))
	}
	return nil
}

// getCommitInfo returns the Commit Info for a given file of the repository
// identified by it's filename
func getCommitInfo(remotePath string, repo *git.Repository) (*OriginFileCommit, error) {

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

	return &OriginFileCommit{
		Author: OriginFileCommitter{
			Name:  returnCommit.Author.Name,
			Email: returnCommit.Author.Email,
		},
		Date: returnCommit.Author.When,
		Hash: returnCommit.Hash.String(),
	}, nil
}

func (file *OriginFile) copyMarkupFile(filesystem billy.Filesystem) error {

	bf, err := filesystem.Open(file.RemotePath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error opening markup file %s", file.RemotePath))
	}

	dirty, err := ioutil.ReadAll(bf)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading markup file %s", file.RemotePath))
	}
	var content string
	contentFormat := file.GetFormat()

	if contentFormat == Markdown {
		content = workarounds.MarkdownPostprocessing(string(dirty))
	} else if contentFormat == Asciidoc {
		content = workarounds.AsciidocPostprocessing(string(dirty))
	}

	content, err = file.ExpandFrontmatter(string(content))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error expanding frontmatter for %s -> %s", file.RemotePath, file.LocalPath))
	}

	err = ioutil.WriteFile(file.LocalPath, []byte(content), standardFilemode)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error writing remote markup file to local file %s -> %s", file.RemotePath, file.LocalPath))
	}
	return nil
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
func (file *OriginFile) ExpandFrontmatter(content string) (expandedFrontmatter string, err error) {

	if file.Commit == nil {
		log.Debug("Git Info is not set, returning without adding it")
		return content, nil
	}

	oldFrontmatter, body, err := splitFrontmatterAndBody(content)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Error expanding front matter"))
	}

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
				file.Commit.Hash,
			),
			// Use lastMod because other variables won't be parsed as date by Hugo
			// Resulting in no date format functions on the file
			file.Commit.Date.Format(time.RFC3339),
			file.Commit.Author.Name,
			file.Commit.Author.Email),
		nil

}

func splitFrontmatterAndBody(content string) (frontmatter string, body string, err error) {
	contentFrontmatter, err := pageparser.ParseFrontMatterAndContent(strings.NewReader(content))
	if err != nil {
		return "", "", errors.Wrap(err, fmt.Sprintf("While splitting frontmatter for content: %s", content))
	}

	// No frontmatter found, return old content
	if contentFrontmatter.Content == nil {
		return "", content, nil
	}

	contentMarshaled, err := yaml.Marshal(contentFrontmatter.FrontMatter)
	if err != nil {
		return "", "", errors.Wrap(err, fmt.Sprintf("Error while marshalling frontmatter to YAML"))
	}

	return string(contentMarshaled), string(contentFrontmatter.Content), nil
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

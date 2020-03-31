package helpers

import (
	"errors"
	"fmt"
	"strings"

	hugo "github.com/gohugoio/hugo/commands"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// FileIsWhitelisted returns true if the filename is in the whitelisted based on its suffix
func FileIsWhitelisted(filename string, whitelist []string) bool {
	for _, whitelisted := range whitelist {
		if strings.HasSuffix(strings.ToLower(filename), strings.ToLower(whitelisted)) {
			return true
		}
	}
	return false
}

// IsMarkdown returns true if a file is a Markdown file
func IsMarkdown(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".md"))
}

// IsAsciidoc returns true if a file is a Asciidoc file
func IsAsciidoc(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".adoc")) ||
		strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".asciidoc")) ||
		strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".asc"))
}

// GetCommitInfo returns the Commit Info for a given file of the repository
// identified by it's filename
func GetCommitInfo(r *git.Repository, filename string) (*object.Commit, error) {

	// TODO what is wrong here?
	cIter, err := r.Log(&git.LogOptions{
		FileName: &filename,
		All:      true,
	})

	if err != nil {
		return nil, fmt.Errorf("Error while opening %s from git log: %s", filename, err)
	}

	var returnCommit *object.Commit

	err = cIter.ForEach(func(commit *object.Commit) error {
		if commit == nil {
			return errors.New("Commit is nil")
		}
		returnCommit = commit
		return nil
	},
	)
	defer cIter.Close()

	if err != nil {
		return nil, err
	}

	return returnCommit, nil
}

// HugoRun runs Hugo like the command line interface
func HugoRun(args []string) error {
	response := hugo.Execute(args)
	return response.Err
}

package helpers

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// ExpandFrontmatter expands the existing frontmatter with the parameters given
func ExpandFrontmatter(content string, g *git.Repository, gitFilepath string, commitinfo *object.Commit) string {

	return fmt.Sprintf(`---
%s

gitRemote : "%s"
gitPath : "%s"
gitLastCommitDate : "%s"
gitLastCommitAuthor : "%s"
gitLastCommitAuthorEmail : "%s"
---

`+content,
		getOldFrontMatter(content),
		"",
		gitFilepath,
		commitinfo.Author.When,
		commitinfo.Author.Name,
		commitinfo.Author.Email)

}

func getOldFrontMatter(content string) string {
	// TODO Convert from toml, yaml, etc
	return `
fakeOldFrontmatter : "FIXME"
`
}

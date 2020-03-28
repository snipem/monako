package helpers

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// ExpandFrontmatter expands the existing frontmatter with the parameters given
func ExpandFrontmatter(content string, commitinfo *object.Commit) string {

	return fmt.Sprintf(`---
title : "test"
creationDate : "%s"
---

`+content, commitinfo.Author.When)

}

package compose

// run: go test ./pkg/compose -run TestGetWebLink

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLocalPath tests if the local file path calculation for remote files is correct
func TestLocalPath(t *testing.T) {

	equalPath(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", ".", ".", "filename.md"),
		"Simple setup, always first level")

	equalPath(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", "docs", ".", "docs/filename.md"),
		"With remote 'docs' folder")

	equalPath(t,
		"/tmp/compose/docs/filename.md",
		getLocalFilePath("/tmp/compose", ".", ".", "docs/filename.md"),
		"With remote 'docs' folder, but keep structure")

	equalPath(t,
		"compose/filename.md",
		getLocalFilePath("./compose", ".", ".", "filename.md"),
		"Path is relative")

	equalPath(t,
		"/tmp/compose/localTarget/filename.md",
		getLocalFilePath("/tmp/compose", ".", "localTarget", "filename.md"),
		"Local Target folder")

	equalPath(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", ".", "", "filename.md"),
		"Empty local target folder")
}

// equalPath is like assert.Equal but with ignoring operation system specifc pathes.
// On Unix "/" and Windows "\" systems this check compares pathes either way.
func equalPath(t *testing.T, expected string, actual string, msg string) {

	assert.Equal(t,
		filepath.ToSlash(expected),
		filepath.ToSlash(actual),
		msg,
	)
}

func TestFrontmatterExpanding(t *testing.T) {

	t.Run("No Frontmatter", func(t *testing.T) {
		content := `=== Body Content
123`
		frontmatter, body := splitFrontmatterAndBody(content)
		assert.Equal(t,
			`=== Body Content
123`,
			body,
			"",
		)

		assert.Empty(t, frontmatter)

	})

	t.Run("Frontmatter with simple file", func(t *testing.T) {
		content := `---
simple: content
content: linetwo
---

=== Body Content
123`
		frontmatter, body := splitFrontmatterAndBody(content)
		assert.Equal(t,
			`
=== Body Content
123`,
			body,
			"",
		)

		assert.Contains(t, frontmatter, "simple: content\n")
		assert.Contains(t, frontmatter, "content: linetwo\n")

	})

	t.Run("Frontmatter garbled file with frontmatter style control signs", func(t *testing.T) {
		content := `---
simple: content
content: linetwo
---

=== Body Content
123
Here be the --- control signs
---
Also on new line`

		frontmatter, body := splitFrontmatterAndBody(content)
		assert.Equal(t,
			`
=== Body Content
123
Here be the --- control signs
---
Also on new line`,
			body,
			"",
		)

		assert.Contains(t, frontmatter, "simple: content\n")
		assert.Contains(t, frontmatter, "content: linetwo\n")

	})

	t.Run("Frontmatter TOML to YAML", func(t *testing.T) {
		content := `+++
simple = "content"
content = "linetwo"
+++

=== Body Content
123
Here be the +++ control signs
+++
Also on new line`

		frontmatter, _ := splitFrontmatterAndBody(content)
		assert.Contains(t, frontmatter, "simple: content\n")
		assert.Contains(t, frontmatter, "content: linetwo\n")
	})
}

func TestGetWebLink(t *testing.T) {

	cases := []struct {
		gitURL, branch, remotePath, Expected string
	}{
		{"https://github.com/snipem/monako-test.git",
			"master", "test_doc_asciidoc.adoc",
			"https://github.com/snipem/monako-test/blob/master/test_doc_asciidoc.adoc"},

		{"https://gitlab.com/snipem/monako-test.git",
			"test-branch", "README.md",
			"https://gitlab.com/snipem/monako-test/blob/test-branch/README.md"},

		{"https://bitbucket.org/snipem/monako-test.git",
			"develop", "README.md",
			"https://bitbucket.org/snipem/monako-test/src/develop/README.md"},

		{"/file/local",
			"develop", "README.md",
			""},

		{"git@github.com:snipem/monako-test.git",
			"master", "README.md",
			""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s, %s, %s -> %s", tc.gitURL, tc.branch, tc.remotePath, tc.Expected), func(t *testing.T) {
			assert.Equal(t, tc.Expected, getWebLinkForFileInGit(tc.gitURL, tc.branch, tc.remotePath))
		})
	}
}

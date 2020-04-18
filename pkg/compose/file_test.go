package compose

// run: MONAKO_TEST_REPO="$HOME/temp/monako-testrepos/monako-test" go test ./pkg/compose -run TestCommitInfo

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
		frontmatter, body, err := splitFrontmatterAndBody(content)
		assert.NoError(t, err)

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
		frontmatter, body, err := splitFrontmatterAndBody(content)
		assert.NoError(t, err)

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

	t.Run("JSON Frontmatter", func(t *testing.T) {
		content := `{
			"categories": [
			   "Development",
			   "Docs"
			],
			"description": "This is the description",
			"date": "2020-04-06",
			"title": "This is the title"
		 }

=== Body Content
123

Inline Json Test {"date": "today"}
Bottom line
`
		frontmatter, body, err := splitFrontmatterAndBody(content)
		assert.NoError(t, err)

		assert.Equal(t,
			`
=== Body Content
123

Inline Json Test {"date": "today"}
Bottom line
`,
			body,
			"",
		)

		assert.Contains(t, frontmatter, "description: This is the description\n")
		assert.Contains(t, frontmatter, "title: This is the title\n")

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

		frontmatter, body, err := splitFrontmatterAndBody(content)
		assert.NoError(t, err)

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

		frontmatter, _, err := splitFrontmatterAndBody(content)
		assert.NoError(t, err)

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

func TestGetCommitWebLink(t *testing.T) {

	cases := []struct {
		gitURL, commitID, Expected string
	}{
		{"https://github.com/snipem/monako-test.git",
			"b744ffe4761cb3a282dcb30ac23b129ec19c9a53",
			"https://github.com/snipem/monako-test/commit/b744ffe4761cb3a282dcb30ac23b129ec19c9a53"},

		{"https://gitlab.com/snipem/monako-test.git",
			"1559b863ff3a9cc1c077ebc480215fd54b621693",
			"https://gitlab.com/snipem/monako-test/commit/1559b863ff3a9cc1c077ebc480215fd54b621693"},

		{"https://bitbucket.org/snipem/monako-test.git",
			"e99f32612df02ee18de15bd42326a10e4195be3d",
			"https://bitbucket.org/snipem/monako-test/commits/e99f32612df02ee18de15bd42326a10e4195be3d"},

		{"/file/local",
			"commitID4711",
			""},

		{"git@github.com:snipem/monako-test.git",
			"commitID4711",
			""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s, %s -> %s", tc.gitURL, tc.commitID, tc.Expected), func(t *testing.T) {
			assert.Equal(t, tc.Expected, getWebLinkForGitCommit(tc.gitURL, tc.commitID))
		})
	}
}

func TestCommitInfo(t *testing.T) {

	testConfig, _ := getTestConfig(t)
	origin := &testConfig.Origins[0]

	_, err := origin.CloneDir()
	assert.NoError(t, err)

	t.Run("Test Commit Info", func(t *testing.T) {
		commit, err := getCommitInfo("README.md", origin.repo)
		assert.NoError(t, err)
		assert.Contains(t, commit.Author.Email, "@")
		assert.NotNil(t, commit.Date)
		assert.NotNil(t, commit.Hash)
		assert.NotNil(t, commit.Author.Name)
	})

	t.Run("Non Existing file", func(t *testing.T) {
		commit, err := getCommitInfo("THIS FILE WILL NEVER EXIST. fake", origin.repo)
		assert.Error(t, err)
		assert.Nil(t, commit)
	})

	t.Run("No repo", func(t *testing.T) {
		commit, err := getCommitInfo("README.md", nil)
		assert.Error(t, err)
		assert.Nil(t, commit)
	})

}

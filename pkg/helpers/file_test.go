package helpers

// run: make test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
)

var g *git.Repository
var fs billy.Filesystem

func TestMain(m *testing.M) {
	// Setup git clone of repo
	setup()
	os.Exit(m.Run())
}

func setup() {
	g, fs = CloneDir("https://github.com/snipem/monako.git", "master", "", "")
}

func TestCloneDir(t *testing.T) {

	root, err := fs.Chroot(".")
	assert.NoError(t, err, "Could not chroot dir")

	files, err := root.ReadDir(".")
	assert.NoError(t, err, "Could not read dir")
	assert.NotZero(t, len(files), "Should have cloned some files")

}

func TestIsMarkdown(t *testing.T) {
	assert.True(t, isMarkdown("markdown.md"), "Check should be true")
	assert.True(t, isMarkdown("markdown.MD"), "Check should be true")
	assert.False(t, isMarkdown("somefolderwith.md-init/somefile.tmp"), "Asciidoc not detected correctly")
}

func TestIsAsciidoc(t *testing.T) {
	assert.True(t, isAsciidoc("asciidoc.adoc"), "Check should be true")
	assert.True(t, isAsciidoc("asciidoc.ADOC"), "Check should be true")
	assert.False(t, isAsciidoc("somefolderwith.adoc-init/somefile.tmp"), "Asciidoc not detected correctly")
}

func TestGitCommiter(t *testing.T) {
	fileName := "README.md"

	ci, err := GetCommitInfo(g, fileName)

	assert.NoError(t, err, "Could not retrieve commit info")
	assert.Contains(t, ci.Committer.Email, "@")

}

func TestGitCommiterFileNotFound(t *testing.T) {
	fileName := "Not existing file...."
	_, err := GetCommitInfo(g, fileName)

	assert.Error(t, err, "Expect error for non existing file")
}

func TestGitCommiterSubfolder(t *testing.T) {
	fileName := "test/config.menu.local.md"
	ci, err := GetCommitInfo(g, fileName)

	assert.NoError(t, err, "Could not retrieve commit info")
	assert.Contains(t, ci.Committer.Email, "@")
}

// TestCopyDir is a test for testing the copying capability of a single directory
func TestCopyDir(t *testing.T) {

	// TODO Get own temporary test folder
	targetDir := filepath.Join(os.TempDir(), "tmp/testrun/", t.Name())
	defer os.RemoveAll(targetDir)

	var whitelist = []string{".md", ".png"}

	t.Run("start in single directory 'test'", func(t *testing.T) {
		CopyDir(g, fs, "test", targetDir, whitelist)
		expectedTargetFile := filepath.Join(targetDir, "test_docs/test_doc_markdown.md")
		b, err := ioutil.ReadFile(expectedTargetFile)

		assert.NoError(t, err, "File not found")
		assert.Contains(t, string(b), "# Markdown Doc 1")
	})

	t.Run("start in deeper directory 'test/test_docs/'", func(t *testing.T) {
		CopyDir(g, fs, "test/test_docs/", targetDir, whitelist)
		expectedTargetFile := filepath.Join(targetDir, "/test_doc_markdown.md")
		b, err := ioutil.ReadFile(expectedTargetFile)

		assert.NoError(t, err, "File not found")
		assert.Contains(t, string(b), "# Markdown Doc 1")
	})

}

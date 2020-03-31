package compose

// run: go test -v ./pkg/compose

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
)

var o *Origin

func TestMain(m *testing.M) {
	// Setup git clone of repo
	setup()
	os.Exit(m.Run())
}

func setup() {
	o = NewOrigin("https://github.com/snipem/monako.git", "master", ".", "test/docs")
	o.CloneDir()
}

// TestCopyDir is a test for testing the copying capability of a single directory
func TestCopyDir(t *testing.T) {
	t.Skip("Skip for now, is hard to understand")

	// TODO Get own temporary test folder
	// targetDir := filepath.Join(os.TempDir(), "tmp/testrun/", t.Name())
	// defer os.RemoveAll(targetDir)

	var whitelist = []string{".md", ".png"}

	t.Run("start in single directory 'test'", func(t *testing.T) {
		o.SourceDir = "test"
		o.FileWhitelist = whitelist
		tempDir := filet.TmpDir(t, "")
		o.ComposeDir(tempDir)
		expectedTargetFile := filepath.Join(tempDir, "compose", "test_docs/test_doc_markdown.md")
		b, err := ioutil.ReadFile(expectedTargetFile)

		assert.NoError(t, err, "File not found")
		assert.Contains(t, string(b), "# Markdown Doc 1")
	})

	t.Run("start in deeper directory 'test/test_docs/'", func(t *testing.T) {
		o.SourceDir = "test/test_docs/"
		o.FileWhitelist = whitelist
		tempDir := filet.TmpDir(t, "")
		o.ComposeDir(tempDir)
		expectedTargetFile := filepath.Join(tempDir, "compose", "/test_doc_markdown.md")
		b, err := ioutil.ReadFile(expectedTargetFile)

		assert.NoError(t, err, "File not found")
		assert.Contains(t, string(b), "# Markdown Doc 1")
	})

}

func TestLocalPath(t *testing.T) {

	assert.Equal(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/", "compose", ".", "filename.md"),
		"Simple setup, always first level")

	assert.Equal(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/", "compose", "docs", "docs/filename.md"),
		"With remote 'docs' folder")

	assert.Equal(t,
		"/tmp/compose/docs/filename.md",
		getLocalFilePath("/tmp/", "compose", ".", "docs/filename.md"),
		"With remote 'docs' folder, but keep structure")

	assert.Equal(t,
		"compose/filename.md",
		getLocalFilePath(".", "compose", ".", "filename.md"),
		"Path is relative")
}

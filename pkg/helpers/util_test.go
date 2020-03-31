package helpers

import (
	"os"
	"testing"

	"github.com/alecthomas/assert"
)

func TestMain(m *testing.M) {
	// Setup git clone of repo
	setup()
	os.Exit(m.Run())
}

func setup() {
}
func TestIsMarkdown(t *testing.T) {
	assert.True(t, IsMarkdown("markdown.md"), "Check should be true")
	assert.True(t, IsMarkdown("markdown.MD"), "Check should be true")
	assert.False(t, IsMarkdown("somefolderwith.md-init/somefile.tmp"), "Asciidoc not detected correctly")
}

func TestIsAsciidoc(t *testing.T) {
	assert.True(t, IsAsciidoc("asciidoc.adoc"), "Check should be true")
	assert.True(t, IsAsciidoc("asciidoc.ADOC"), "Check should be true")
	assert.False(t, IsAsciidoc("somefolderwith.adoc-init/somefile.tmp"), "Asciidoc not detected correctly")
}

// TODO Use direct repository

// func TestGitCommiter(t *testing.T) {
// 	fileName := "README.md"

// 	ci, err := GetCommitInfo(o.repo, fileName)

// 	assert.NoError(t, err, "Could not retrieve commit info")
// 	assert.Contains(t, ci.Committer.Email, "@")

// }

// func TestGitCommiterFileNotFound(t *testing.T) {
// 	fileName := "Not existing file...."
// 	_, err := GetCommitInfo(o.repo, fileName)

// 	assert.Error(t, err, "Expect error for non existing file")
// }

// func TestGitCommiterSubfolder(t *testing.T) {
// 	fileName := "test/config.menu.local.md"
// 	ci, err := GetCommitInfo(o.repo, fileName)

// 	assert.NoError(t, err, "Could not retrieve commit info")
// 	assert.Contains(t, ci.Committer.Email, "@")
// }

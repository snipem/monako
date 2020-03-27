package helpers

// run: make test

import (
	"os"
	"strings"
	"testing"

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

	if err != nil {
		t.Error(err)
	}

	files, err := root.ReadDir(".")

	if err != nil {
		if len(files) == 0 {
			t.Errorf("No files checked out")
		}
	}
}

func TestIsMarkdown(t *testing.T) {
	if !isMarkdown("markdown.md") || !isMarkdown("markdown.MD") {
		t.Error("Markdown not detected correctly")
	}

	if isMarkdown("somefolderwith.md-init/somefile.tmp") {
		t.Error("Markdown not detected correctly")
	}
}

func TestIsAsciidoc(t *testing.T) {
	if !isAsciidoc("asciidoc.adoc") || !isAsciidoc("example.ADOC") {
		t.Error("Asciidoc not detected correctly")
	}

	if isAsciidoc("somefolderwith.adoc-init/somefile.tmp") {
		t.Error("Asciidoc not detected correctly")
	}
}

func TestGitCommiter(t *testing.T) {
	fileName := "README.md"

	ci, err := GetCommitInfo(g, fileName)

	if err != nil {
		t.Error(err)
	}

	mail := ci.Committer.Email
	if !strings.Contains(mail, "@") {
		t.Errorf("Commiter %s does not contain @", mail)
	}

}

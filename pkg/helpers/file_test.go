package helpers

// run: make test

import (
	"fmt"
	"io/ioutil"
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

func TestGitCommiterFileNotFound(t *testing.T) {
	fileName := "Not existing file...."

	_, err := GetCommitInfo(g, fileName)

	if err == nil {
		t.Error("No error given")
	}

	if strings.Contains(err.Error(), "EOF") {
		t.Error("Err contains EOF and is too technical")
	}

}

func TestGitCommiterSubfolder(t *testing.T) {
	fileName := ".github/workflows/main.yml"

	ci, err := GetCommitInfo(g, fileName)

	if err != nil {
		t.Error(err)
	}

	mail := ci.Committer.Email
	if !strings.Contains(mail, "@") {
		t.Errorf("Commiter %s does not contain @", mail)

	}
}

func TestCopyDir(t *testing.T) {

	target := "tmp/testrun/"
	var whitelist = []string{".md", ".png"}
	CopyDir(g, fs, "test", target, whitelist)

	_, err := ioutil.ReadFile(target + "test_docs/test_doc_markdown.md")
	if err != nil {
		fmt.Print(err)
		t.Errorf("To be copied file not found")
	}

}

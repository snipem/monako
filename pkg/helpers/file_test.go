package helpers

// run: make test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
)

var g *git.Repository
var fs billy.Filesystem

func TestMain(m *testing.M) {
	// Setup git clone of repo
	// log.SetReportCaller(true)
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
	// fileName := "test/config.menu.local.md"
	fileName := ".github/workflows/main.yml"

	ci, err := GetCommitInfo(g, fileName)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	mail := ci.Committer.Email
	if !strings.Contains(mail, "@") {
		t.Errorf("Commiter %s does not contain @", mail)

	}
}

// TestCopyDir is a test for testing the copying capability of a single directory
func TestCopyDir(t *testing.T) {

	targetDir := filepath.Join(os.TempDir(), "tmp/testrun/", t.Name())
	defer os.RemoveAll(targetDir)

	var whitelist = []string{".md", ".png"}
	CopyDir(g, fs, "test", targetDir, whitelist)

	expectedTargetFile := filepath.Join(targetDir, "test_docs/test_doc_markdown.md")

	b, err := ioutil.ReadFile(expectedTargetFile)
	if err != nil {
		t.Errorf("Expected file %s not found", expectedTargetFile)
		t.FailNow()
	}

	if !strings.Contains(string(b), "# Markdown Doc 1") {
		t.Errorf("Wrong file copied under right name")
	}
}

// TestCopyDirWithSubfolderSource tests if a source can be a deeper path like "test/test_docs" instead of
// starting at highest level or just one directory deep
func TestCopyDirWithSubfolderSource(t *testing.T) {
	targetDir := filepath.Join(os.TempDir(), "tmp/testrun/", t.Name())
	defer os.RemoveAll(targetDir)

	var whitelist = []string{".md", ".png"}
	CopyDir(g, fs, "test/test_docs/", targetDir, whitelist)

	expectedTargetFile := filepath.Join(targetDir, "/test_doc_markdown.md")

	b, err := ioutil.ReadFile(expectedTargetFile)
	if err != nil {
		t.Errorf("Expected file %s not found", expectedTargetFile)
		t.FailNow()
	}

	if !strings.Contains(string(b), "# Markdown Doc 1") {
		t.Errorf("Wrong file copied under right name")
	}
}

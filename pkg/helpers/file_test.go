package helpers

// run: make test

import (
	"testing"
)

func TestCloneDir(t *testing.T) {

	// Use current dir as test repo, navigate two folders up because we are in /pkg/helpers
	_, fs := CloneDir("https://github.com/snipem/monako.git", "master", "", "")

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

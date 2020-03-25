package helpers

// run: make test

import (
	"testing"
)

func TestCloneDir(t *testing.T) {

	// Use current dir as test repo, navigate two folders up because we are in /pkg/helpers
	git, fs := CloneDir("https://github.com/snipem/monako.git", "master", "", "")

	if git != nil {
		t.Log("Git Clone returned")
	}

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

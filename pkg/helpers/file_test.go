package helpers

// run: make test

import (
	"os"
	"testing"
)

func TestCloneDir(t *testing.T) {

	wd, err := os.Getwd()

	// Use current dir as test repo, navigate two folders up because we are in /pkg/helpers
	git, fs := CloneDir("file://"+wd+"/../..", "master", "", "")

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

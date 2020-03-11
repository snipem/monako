package main

import "testing"

func TestMarkdownFix(t *testing.T) {

	noNeedToClean := "!(caption example)[http://url/image.jpg]"
	clean := string(MarkdownPostprocessing([]byte(noNeedToClean)))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

	dirty := "!(caption example)[lokalfolderurl/image.png]"
	clean = string(MarkdownPostprocessing([]byte(dirty)))
	want := "!(caption example)[lokalfolderurl/image.png]"

	if clean != want {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
	}

}

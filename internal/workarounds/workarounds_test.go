package workarounds

// run: make test

import "testing"

func TestAsciiDocImageFix(t *testing.T) {

	noNeedToClean := "image:http://url/image.jpg[image,width=634,height=346]"
	needToClean := "image:image2.png[image,width=634,height=346]"

	clean := string(AsciidocPostprocessing([]byte(noNeedToClean)))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

	clean = string(AsciidocPostprocessing([]byte(needToClean)))
	want := "image:../image2.png[image,width=634,height=346]"

	if clean != want {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
	}
}

func TestMarkdownFix(t *testing.T) {

	noNeedToClean := "!(caption example)[http://url/image.jpg]"
	clean := string(MarkdownPostprocessing([]byte(noNeedToClean)))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

}
func TestMarkdownFixDontDo(t *testing.T) {

	dirty := "!(caption example)[lokalfolderurl/image.png]"
	want := "!(caption example)[../lokalfolderurl/image.png]"

	clean := string(MarkdownPostprocessing([]byte(dirty)))

	if clean != want {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
	}

}
func TestMarkdownFixInnerLinks(t *testing.T) {
	// Dont fix local links

	noNeedToClean := "[[1]](#1)"

	clean := string(MarkdownPostprocessing([]byte(noNeedToClean)))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}
}

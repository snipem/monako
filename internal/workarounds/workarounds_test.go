package workarounds

// run: make test

import "testing"

func TestAsciiDocImageFix(t *testing.T) {

	noNeedToClean := "image:http://url/image.jpg[image,width=634,height=346]"
	needToClean := "image:image2.png[image,width=634,height=346]"

	clean := string(asciidocPostprocessing([]byte(noNeedToClean)))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

	clean = string(asciidocPostprocessing([]byte(needToClean)))
	want := "image:../image2.png[image,width=634,height=346]"

	if clean != want {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
	}
}

func TestMarkdownFix(t *testing.T) {

	noNeedToClean := "!(caption example)[http://url/image.jpg]"
	clean := string(markdownPostprocessing([]byte(noNeedToClean)))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

	dirty := "!(caption example)[lokalfolderurl/image.png]"
	clean = string(markdownPostprocessing([]byte(dirty)))
	want := "!(caption example)[lokalfolderurl/image.png]"

	if clean != want {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
	}

}

func TestAddAsciidocFix(t *testing.T) {
	t.Skip("Skip test since there is no actual check")
	// TODO add actual test
	addFixForADocTocToTheme()
}

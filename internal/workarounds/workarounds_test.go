package workarounds

// run: make test

import (
	"testing"
)

func TestAsciiDocImageFix(t *testing.T) {

	noNeedToClean := "image:http://url/image.jpg[image,width=634,height=346]"
	needToClean := "image:image2.png[image,width=634,height=346]"

	clean := AsciidocPostprocessing(noNeedToClean)

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

	clean = AsciidocPostprocessing(needToClean)
	want := "image:../image2.png[image,width=634,height=346]"

	if clean != want {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
	}

	t.Run("Single colon and local syntax local dir", func(t *testing.T) {

		needToClean := "image:./image2.png[image,width=634,height=346]"

		clean = AsciidocPostprocessing(needToClean)
		want := "image:../image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})

	t.Run("Double colons", func(t *testing.T) {

		needToClean := "image::image2.png[image,width=634,height=346]"

		clean = AsciidocPostprocessing(needToClean)
		want := "image::../image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})

	t.Run("Double colons and local syntax local dir", func(t *testing.T) {

		needToClean := "image::./image2.png[image,width=634,height=346]"

		clean = AsciidocPostprocessing(needToClean)
		want := "image::../image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})

	t.Run("Double colons with absolute url, do nothing", func(t *testing.T) {

		needToClean := "image::http://absolute/url/image2.png[image,width=634,height=346]"

		clean = AsciidocPostprocessing(needToClean)
		want := "image::http://absolute/url/image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})
}

func TestMarkdownFix(t *testing.T) {

	noNeedToClean := "![caption example](http://url/image.jpg)"
	clean := MarkdownPostprocessing(noNeedToClean)

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

}
func TestMarkdownFixDontDo(t *testing.T) {

	dirty := "![caption example](lokalfolderurl/image.png)"
	want := "![caption example](../lokalfolderurl/image.png)"

	clean := string(MarkdownPostprocessing(dirty))

	if clean != want {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
	}

}
func TestMarkdownFixInnerLinks(t *testing.T) {
	// Dont fix local links

	t.Run("Pseudo references", func(t *testing.T) {
		noNeedToClean := "[[1]](#1)"

		clean := string(MarkdownPostprocessing(noNeedToClean))

		if clean != noNeedToClean {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
		}

	})

	t.Run("Alpbabet links", func(t *testing.T) {

		noNeedToClean := "[details](#teamname)"
		clean := string(MarkdownPostprocessing(noNeedToClean))
		if clean != noNeedToClean {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
		}

	})
}

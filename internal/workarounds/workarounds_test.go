package workarounds

// run: make test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	t.Run("Single colon and local syntax local dir", func(t *testing.T) {

		needToClean := "image:./image2.png[image,width=634,height=346]"

		clean = string(AsciidocPostprocessing([]byte(needToClean)))
		want := "image:../image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})

	t.Run("Double colons", func(t *testing.T) {

		needToClean := "image::image2.png[image,width=634,height=346]"

		clean = string(AsciidocPostprocessing([]byte(needToClean)))
		want := "image::../image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})

	t.Run("Double colons and local syntax local dir", func(t *testing.T) {

		needToClean := "image::./image2.png[image,width=634,height=346]"

		clean = string(AsciidocPostprocessing([]byte(needToClean)))
		want := "image::../image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})

	t.Run("Double colons with absolute url, do nothing", func(t *testing.T) {

		needToClean := "image::http://absolute/url/image2.png[image,width=634,height=346]"

		clean = string(AsciidocPostprocessing([]byte(needToClean)))
		want := "image::http://absolute/url/image2.png[image,width=634,height=346]"

		if clean != want {
			t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, want)
		}

	})
}

func TestMarkdownFix(t *testing.T) {

	noNeedToClean := "![caption example](http://url/image.jpg)"
	clean := string(MarkdownPostprocessing([]byte(noNeedToClean)))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}

}
func TestMarkdownFixDontDo(t *testing.T) {

	dirty := "![caption example](lokalfolderurl/image.png)"
	want := "![caption example](../lokalfolderurl/image.png)"

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

func TestFakeAsciidoctorBin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Can't test this on Windows")
	}
	_, isGithubWorkflow := os.LookupEnv("GITHUB_WORKFLOW")
	if isGithubWorkflow {
		t.Skip("Don't run this test on Github Actions")
	}

	fakePath := AddFakeAsciidoctorBinForDiagramsToPath("http://complexbasepath/path/bla")

	assert.FileExists(t, fakePath)

	// This is how Hugo does it: https://github.com/gohugoio/hugo/blob/master/markup/asciidoc/convert.go#L90
	// We wont to trick Hugo into the same behaviour
	resolvedPath, err := exec.LookPath("asciidoctor")
	assert.NoError(t, err)
	assert.Equal(t, fakePath, resolvedPath)

	read, err := ioutil.ReadFile(fakePath)
	assert.NoError(t, err)
	assert.Contains(t, string(read), "\\/path\\/bla")

}

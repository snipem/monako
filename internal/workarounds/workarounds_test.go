package workarounds

// run: make test

import (
	"io/ioutil"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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

	noNeedToClean := "[[1]](#1)"

	clean := string(MarkdownPostprocessing(noNeedToClean))

	if clean != noNeedToClean {
		t.Errorf("Clean was incorrect, got: %s, want: %s.", clean, noNeedToClean)
	}
}

func TestFakeAsciidoctorBin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Can't test this on Windows")
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

// func TestRenderMermaidLikeGitlab(t *testing.T) {

// 	before = "```" + `mermaid
// graph TD;
//     A-->B;
//     A-->C;
//     B-->D;
//     C-->D;
// 	` + "```" + `

// 	after = `
// 	{{< mermaid >}}
// 	{{< /mermaid >}}
// `

// }

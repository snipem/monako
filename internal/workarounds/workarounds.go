package workarounds

// run: make test

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// AsciidocPostprocessing fixes common errors with Hugo processing vanilla Asciidoc
// 1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
// for pretty urls
func AsciidocPostprocessing(dirty []byte) []byte {

	var d = string(dirty)

	// TODO: These workarounds will be removed when the asciidoc fix for pathes arrives in Hugo
	// https://github.com/gohugoio/hugo/pull/6561
	// This workaround is still needed as long as we're relying on the asciidoctor diagram fix
	// since the diagram fix does not work with "ugly" urls

	// Really quick and dirty. There is a problem with Go regexp look ahead
	// Preserve
	d = strings.ReplaceAll(d, "image::http", "image+______http")
	d = strings.ReplaceAll(d, "image:http", "image_______http")

	// Replace
	d = strings.ReplaceAll(d, "image:", "image:../")

	// Restore
	d = strings.ReplaceAll(d, "image_______http", "image:http")
	d = strings.ReplaceAll(d, "image+______http", "image::http")

	// Fix for colons being moved
	d = strings.ReplaceAll(d, "image:../:", "image::../")

	// Fix for ./ syntax. It's getting uglier
	d = strings.ReplaceAll(d, "image::.././", "image::../")
	d = strings.ReplaceAll(d, "image:.././", "image:../")

	return []byte(d)
}

// MarkdownPostprocessing fixes common errors with Hugo processing vanilla Markdown
//  1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
// for pretty urls
func MarkdownPostprocessing(dirty []byte) []byte {
	var d = string(dirty)

	// FIXME really quick and dirty. There is a problem with Go regexp look ahead
	d = strings.ReplaceAll(d, "](http", ")_______http")
	d = strings.ReplaceAll(d, "]](", "]]_______(")

	d = strings.ReplaceAll(d, "](", "](../")

	d = strings.ReplaceAll(d, ")_______http", "](http")
	d = strings.ReplaceAll(d, "]]_______(", "]](")

	return []byte(d)
}

// AddFakeAsciidoctorBinForDiagramsToPath adds a fake asciidoctor bin to the PATH
// to trick Hugo into taking this one. This makes it possible to manipulate the parameters
// for asciidoctor while being called from Hugo.
func AddFakeAsciidoctorBinForDiagramsToPath(baseURL string) string {

	if runtime.GOOS == "windows" {
		log.Println("Can't apply asciidoctor diagram workaround on Windows")
		return ""
	}

	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	path := url.Path

	// Single slashes will add up to "//" which some webservers don't support
	if path == "/" {
		path = ""
	}
	escapedPath := strings.ReplaceAll(path, "/", "\\/")
	escapedPath = strings.ReplaceAll(escapedPath, "\"", "\\\"")

	originalAsciidoctorPath, _ := exec.LookPath("asciidoctor")
	// Ignore error otherwise runs with no asciidoctor installed will fail
	if err != nil {
		log.Debugf("No asciidoctor found in path")
		return ""
	}

	// Asciidoctor attributes: https://asciidoctor.org/docs/user-manual/#builtin-attributes

	// TODO: Use variables
	shellscript := fmt.Sprintf(`#!/bin/bash
	# inspired by: https://zipproth.de/cheat-sheets/hugo-asciidoctor/#_how_to_make_hugo_use_asciidoctor_with_extensions
	set -e

	# Use empty css to trick asciidoctor into using none without error
	echo "" > empty.css

	# This trick only works with the relative dir workarounds
	%s -B . \
		-r asciidoctor-diagram \
		-a nofooter \
		-a stylesheet=empty.css \
		--safe \
		--trace \
		- | sed -E -e "s/img src=\"([^/]+)\"/img src=\"%s\/diagram\/\1\"/"

	# For some reason static is not parsed with integrated Hugo
	mkdir -p compose/public/diagram
	
	# Hopefully this will also be fixed by https://github.com/gohugoio/hugo/pull/6561
	if ls *.svg >/dev/null 2>&1; then
	  mv -f *.svg compose/public/diagram
	fi
	
	if ls *.png >/dev/null 2>&1; then
	  mv -f *.png compose/public/diagram
	fi
	`, originalAsciidoctorPath, escapedPath)

	tempDir := filepath.Join(os.TempDir(), "monako_asciidoctor_fake_binary")
	err = os.Mkdir(tempDir, os.FileMode(0700))
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Error creating asciidoctor fake dir : %s", err)
	}
	fakeBinary := filepath.Join(tempDir, "asciidoctor")

	err = ioutil.WriteFile(fakeBinary, []byte(shellscript), os.FileMode(0700))
	if err != nil {
		log.Fatalf("Error creating asciidoctor fake binary: %s", err)
	}

	os.Setenv("PATH", tempDir+":"+os.Getenv("PATH"))

	log.Debugf("Added temporary binary %s to PATH %s", fakeBinary, os.Getenv("PATH"))

	return fakeBinary

}

package workarounds

// run: make test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

// AsciidocPostprocessing fixes common errors with Hugo processing vanilla Asciidoc
// 1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
// for pretty urls
func AsciidocPostprocessing(dirty []byte) []byte {

	var d = string(dirty)

	// FIXME really quick and dirty. There is a problem with Go regexp look ahead
	d = strings.ReplaceAll(d, "image:http", "image_______http")
	d = strings.ReplaceAll(d, "image:", "image:../")
	d = strings.ReplaceAll(d, "image_______http", "image:http")
	return []byte(d)
}

// MarkdownPostprocessing fixes common errors with Hugo processing vanilla Markdown
//  1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
// for pretty urls
func MarkdownPostprocessing(dirty []byte) []byte {
	var d = string(dirty)

	// FIXME really quick and dirty. There is a problem with Go regexp look ahead
	d = strings.ReplaceAll(d, "](http", ")_______http")
	d = strings.ReplaceAll(d, "](", "](../")
	d = strings.ReplaceAll(d, ")_______http", "](http")

	return []byte(d)
}

func AddFakeAsciidoctorBinForDiagramsToPath(baseURL string) string {

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

	// Asciidoctor attributes: https://asciidoctor.org/docs/user-manual/#builtin-attributes
	// TODO: Are these attributes reasonable?

	shellscript := fmt.Sprintf(`#!/bin/bash
	# inspired by: https://zipproth.de/cheat-sheets/hugo-asciidoctor/#_how_to_make_hugo_use_asciidoctor_with_extensions
	if [ -f /usr/local/bin/asciidoctor ]; then
	  ad="/usr/local/bin/asciidoctor"
	else
	  ad="/usr/bin/asciidoctor"
	fi
	
	$ad -v -B . \
		-r asciidoctor-diagram \
		--no-header-footer \
		--safe \
		--trace \
		-a icons=font \
		-a docinfo=shared \
		-a sectanchors \
		-a experimental=true \
		-a figure-caption! \
		-a source-highlighter=highlightjs \
		-a toc-title! \
		-a stem=mathjax \
		- | sed -E -e "s/img src=\"([^/]+)\"/img src=\"%s\/diagram\/\1\"/"
	
	# For some reason static is not parsed with integrated Hugo
	mkdir -p compose/public/diagram
	
	if ls *.svg >/dev/null 2>&1; then
	  mv -f *.svg compose/public/diagram
	fi
	
	if ls *.png >/dev/null 2>&1; then
	  mv -f *.png compose/public/diagram
	fi
	`, escapedPath)

	tempDir := os.TempDir() + "/asciidoctor_fake_binary"
	err = os.Mkdir(tempDir, os.FileMode(0700))
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Error creating asciidoctor fake dir : %s", err)
	}
	fakeBinary := tempDir + "/asciidoctor"

	err = ioutil.WriteFile(fakeBinary, []byte(shellscript), os.FileMode(0700))
	if err != nil {
		log.Fatalf("Error creating asciidoctor fake binary: %s", err)
	}

	os.Setenv("PATH", tempDir+":"+os.Getenv("PATH"))

	log.Printf("Added temporary binary %s to PATH %s", fakeBinary, os.Getenv("PATH"))

	return fakeBinary

}

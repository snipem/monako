package workarounds

// run: make test

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// AsciidocPostprocessing fixes common errors with Hugo processing vanilla Asciidoc
// 1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
// for pretty urls
func AsciidocPostprocessing(dirty string) (clean string) {

	// DEPCRECATED: These workarounds will be removed when the asciidoc fix for pathes arrives in Hugo
	// https://github.com/gohugoio/hugo/pull/6561
	// This workaround is still needed as long as we're relying on the asciidoctor diagram fix
	// since the diagram fix does not work with "ugly" urls

	// Really quick and dirty. There is a problem with Go regexp look ahead
	// Preserve
	dirty = strings.ReplaceAll(dirty, "image::http", "image+______http")
	dirty = strings.ReplaceAll(dirty, "image:http", "image_______http")

	// Replace
	dirty = strings.ReplaceAll(dirty, "image:", "image:../")

	// Restore
	dirty = strings.ReplaceAll(dirty, "image_______http", "image:http")
	dirty = strings.ReplaceAll(dirty, "image+______http", "image::http")

	// Fix for colons being moved
	dirty = strings.ReplaceAll(dirty, "image:../:", "image::../")

	// Fix for ./ syntax. It's getting uglier
	dirty = strings.ReplaceAll(dirty, "image::.././", "image::../")
	dirty = strings.ReplaceAll(dirty, "image:.././", "image:../")

	clean = dirty

	return clean
}

// MarkdownPostprocessing fixes common errors with Hugo processing vanilla Markdown
//  1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
// for pretty urls
func MarkdownPostprocessing(dirty string) (clean string) {

	// FIXME really quick and dirty. There is a problem with Go regexp look ahead

	// DEPCRECATED: These workarounds will be removed when the asciidoc fix for pathes arrives in Hugo
	// https://github.com/gohugoio/hugo/pull/6561

	// Preserve absolute urls
	dirty = strings.ReplaceAll(dirty, "](http", ")_______http")
	// Preserve anchors
	dirty = strings.ReplaceAll(dirty, "](#", "]_______(#")

	dirty = strings.ReplaceAll(dirty, "](", "](../")

	// Restore absolute urls
	dirty = strings.ReplaceAll(dirty, ")_______http", "](http")
	// Restore anchors
	dirty = strings.ReplaceAll(dirty, "]_______(#", "](#")

	clean = dirty

	return clean
}

// AddFakeAsciidoctorBinForDiagramsToPath adds a fake asciidoctor bin to the PATH
// to trick Hugo into taking this one. This makes it possible to manipulate the parameters
// for asciidoctor while being called from Hugo.
func AddFakeAsciidoctorBinForDiagramsToPath(baseURL string) (fakeBinaryPath string, err error) {

	if runtime.GOOS == "windows" {
		log.Warn("Can't apply asciidoctor diagram workaround on Windows")
		return "", nil
	}

	url, err := url.Parse(baseURL)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Error parsing base url %s", baseURL))
	}
	path := url.Path

	// Single slashes will add up to "//" which some webservers don't support
	if path == "/" {
		path = ""
	}
	escapedPath := strings.ReplaceAll(path, "/", "\\/")
	escapedPath = strings.ReplaceAll(escapedPath, "\"", "\\\"")

	// Asciidoctor attributes: https://asciidoctor.org/docs/user-manual/#builtin-attributes

	shellscript := fmt.Sprintf(`#!/bin/bash
	# inspired by: https://zipproth.de/cheat-sheets/hugo-asciidoctor/#_how_to_make_hugo_use_asciidoctor_with_extensions
	set -e

	# Use first non fake-binary in path as asciidoctorbin
	ad=$(which -a asciidoctor | grep -v monako_asciidoctor_fake_binary | head -n 1)

	# Use empty css to trick asciidoctor into using none without error
	echo "" > empty.css

	# This trick only works with the relative dir workarounds
	$ad -B . \
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
	`, escapedPath)

	tempDir := filepath.Join(os.TempDir(), "monako_asciidoctor_fake_binary")
	err = os.Mkdir(tempDir, os.FileMode(0700))
	if err != nil && !os.IsExist(err) {
		return "", errors.Wrap(err, fmt.Sprintf("Error creating asciidoctor fake dir : %s", tempDir))
	}
	fakeBinaryPath = filepath.Join(tempDir, "asciidoctor")

	err = ioutil.WriteFile(fakeBinaryPath, []byte(shellscript), os.FileMode(0700))
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Error creating asciidoctor fake binary: %s", fakeBinaryPath))
	}

	os.Setenv("PATH", tempDir+":"+os.Getenv("PATH"))

	log.Debugf("Added temporary binary %s to PATH %s", fakeBinaryPath, os.Getenv("PATH"))

	return fakeBinaryPath, nil

}

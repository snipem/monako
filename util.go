package main

import (
	"os"
	"strings"

	"github.com/gohugoio/hugo/commands"
)

func CleanUp() {
	os.RemoveAll("compose")
}

func AsciidocPostprocessing(dirty []byte) []byte {
	return dirty
}

func HugoRun(args []string) {
	// args := []string{"--contentDir", "compose"}
	commands.Execute(args)
}

// MarkdownPostprocessing fixes common errors with Hugo processing vanilla Markdown
//  1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
func MarkdownPostprocessing(dirty []byte) []byte {
	var d = string(dirty)

	// FIXME really quick and dirty, just for testing
	d = strings.ReplaceAll(d, "](http", "]xxxxxxhttp")
	d = strings.ReplaceAll(d, "](", "](../")
	d = strings.ReplaceAll(d, "]xxxxxxhttp", "](http")

	return []byte(d)
}

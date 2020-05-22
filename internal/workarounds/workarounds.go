package workarounds

// run: make test

import (
	"strings"
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

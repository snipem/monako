package main

import (
	"os"
	"strings"
)

func AsciidocPostprocessing(dirty []byte) []byte {

	var d = string(dirty)

	// FIXME really quick and dirty, just for testing
	d = strings.ReplaceAll(d, "image:http", "imagexxxxxxhttp")
	d = strings.ReplaceAll(d, "image:", "image:../")
	d = strings.ReplaceAll(d, "imagexxxxxxhttp", "image:http")
	return []byte(d)
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

func AddFixForAsciiDocTocToTheme() {
	themefile := "compose/themes/hugo-book-6/layouts/partials/docs/html-head.html"
	javascriptFix := `
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
	<script type="text/javascript">
	function fixAsciidocToc() {
		var asciidoc_list = $("#toc > ul");

		if (asciidoc_list.length > 0) {
			// Add new sub item to right nav
			var right_toc = $("body > main > aside.book-toc")[0];
			right_toc.innerHTML += "<nav id='TableOfContents'> </nav>";
		
			// Take content from central asciidoc nav to right side
			var new_nav = $("body > main > aside.book-toc > nav")[0];
			new_nav.append(asciidoc_list[0]);
		
			// Remove all "Table of Contents" from central asciidoc
			$("#toctitle")[0].remove();
		}
	}
	window.onload = fixAsciidocToc;
	</script>
	`

	f, err := os.OpenFile(themefile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(javascriptFix); err != nil {
		panic(err)
	}

}

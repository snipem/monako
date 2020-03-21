package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func asciidocPostprocessing(dirty []byte) []byte {

	var d = string(dirty)

	// FIXME really quick and dirty, just for testing
	d = strings.ReplaceAll(d, "image:http", "imagexxxxxxhttp")
	d = strings.ReplaceAll(d, "image:", "image:../")
	d = strings.ReplaceAll(d, "imagexxxxxxhttp", "image:http")
	return []byte(d)
}

// MarkdownPostprocessing fixes common errors with Hugo processing vanilla Markdown
//  1. Add one level to relative picture img/ -> ../img/ since Hugo adds subfolders
func markdownPostprocessing(dirty []byte) []byte {
	var d = string(dirty)

	// FIXME really quick and dirty, just for testing
	d = strings.ReplaceAll(d, "](http", "]xxxxxxhttp")
	d = strings.ReplaceAll(d, "](", "](../")
	d = strings.ReplaceAll(d, "]xxxxxxhttp", "](http")

	return []byte(d)
}

func addFakeAsciidoctorBinForDiagramsToPath() {

	shellscript := `#!/bin/bash
	# inspired by: https://zipproth.de/cheat-sheets/hugo-asciidoctor/#_how_to_make_hugo_use_asciidoctor_with_extensions
	if [ -f /usr/local/bin/asciidoctor ]; then
	  ad="/usr/local/bin/asciidoctor"
	else
	  ad="/usr/bin/asciidoctor"
	fi
	
	# Use stylesheet=none to prevent asciidoctor from inserting an own stylesheet
	$ad -v -B . -r asciidoctor-diagram -a stylesheet=none -a icons=font -a docinfo=shared -a nofooter -a sectanchors -a experimental=true -a figure-caption! -a source-highlighter=highlightjs -a toc-title! -a stem=mathjax - | sed -E -e "s/img src=\"([^/]+)\"/img src=\"\/diagram\/\1\"/"
	
	# For some reason static is not parsed with monako
	mkdir -p compose/public/diagram
	
	if ls *.svg >/dev/null 2>&1; then
	  mv -f *.svg compose/public/diagram
	fi
	
	if ls *.png >/dev/null 2>&1; then
	  mv -f *.png compose/public/diagram
	fi
	`
	tempDir := os.TempDir() + "/asciidoctor_fake_binary"
	err := os.Mkdir(tempDir, os.FileMode(0700))
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Error creating asciidoctor fake dir : %s", err)
	}
	fakeBinary := tempDir + "/asciidoctor"

	err = ioutil.WriteFile(fakeBinary, []byte(shellscript), os.FileMode(0700))
	if err != nil {
		log.Fatalf("Error creating asciidoctor fake binary: %s", err)
	}
	// TODO Remove file afterwards

	os.Setenv("PATH", tempDir+":"+os.Getenv("PATH"))

	log.Printf("Added temporary binary %s to PATH %s", fakeBinary, os.Getenv("PATH"))

}

func addFixForADocTocToTheme() {
	themefile := "compose/themes/monako-book-6s/layouts/partials/docs/html-head.html"
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

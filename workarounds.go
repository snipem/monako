package main

import (
	"os"
)

// TODO Move to theme
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

package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup git clone of repo
	setup()
	os.Exit(m.Run())
}

func setup() {
}
func TestIsMarkdown(t *testing.T) {
	assert.True(t, IsMarkdown("markdown.md"), "Check should be true")
	assert.True(t, IsMarkdown("markdown.MD"), "Check should be true")
	assert.False(t, IsMarkdown("somefolderwith.md-init/somefile.tmp"), "Asciidoc not detected correctly")
}

func TestIsAsciidoc(t *testing.T) {
	assert.True(t, IsAsciidoc("asciidoc.adoc"), "Check should be true")
	assert.True(t, IsAsciidoc("asciidoc.ADOC"), "Check should be true")
	assert.False(t, IsAsciidoc("somefolderwith.adoc-init/somefile.tmp"), "Asciidoc not detected correctly")
}

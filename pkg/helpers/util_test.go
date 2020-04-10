package helpers

// run: go test ./pkg/helpers/

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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

func TestFileIsWhiteOrBlacklisted(t *testing.T) {
	assert.True(t, FileIsListed("filename.txt", []string{"txt"}))
	assert.False(t, FileIsListed("filename.sh", []string{"txt"}))

	assert.True(t, FileIsListed("filename.md", []string{"md", "adoc"}))
	assert.True(t, FileIsListed("filename.adoc", []string{"md", "adoc"}))
	assert.False(t, FileIsListed("filename.adoc", []string{}))
	assert.False(t, FileIsListed("filename.adoc", []string{"md"}))
}

func TestHugoRun(t *testing.T) {
	assert.NoError(t, HugoRun([]string{"version"}))
	assert.Error(t, HugoRun([]string{"unknown-flag-by-monako-test-case"}))
}

func TestTrace(t *testing.T) {
	assert.NotEqual(t, logrus.GetLevel(), logrus.DebugLevel)
	Trace()
	assert.Equal(t, logrus.GetLevel(), logrus.DebugLevel)
}

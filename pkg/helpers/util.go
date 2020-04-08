package helpers

import (
	"strings"

	hugo "github.com/gohugoio/hugo/commands"
	"github.com/sirupsen/logrus"
)

// FileIsWhitelisted returns true if the filename is in the whitelisted based on its suffix
func FileIsWhitelisted(filename string, whitelist []string) bool {
	for _, whitelisted := range whitelist {
		if strings.HasSuffix(strings.ToLower(filename), strings.ToLower(whitelisted)) {
			return true
		}
	}
	return false
}

// IsMarkdown returns true if a file is a Markdown file
func IsMarkdown(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".md"))
}

// IsAsciidoc returns true if a file is a Asciidoc file
func IsAsciidoc(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".adoc")) ||
		strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".asciidoc")) ||
		strings.HasSuffix(strings.ToLower(filename), strings.ToLower(".asc"))
}

// HugoRun runs Hugo like the command line interface
func HugoRun(args []string) error {
	response := hugo.Execute(args)
	return response.Err
}

// Trace sets trace mode
func Trace() {
	logrus.SetLevel(logrus.DebugLevel)
	// Add line and filename to log
	// This is slow
	logrus.SetReportCaller(true)
}

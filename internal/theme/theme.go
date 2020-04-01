package theme

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/c4milo/unpackit"
	"github.com/snipem/monako/pkg/compose"
)

const monakoMenuDirectory = "monako_menu_directory"
const themeName = "monako-book-master"

// CreateHugoPage extracts the Monako theme and copies the hugoconfig and menuconfig to the needed files
func CreateHugoPage(composeConfig compose.Config, menuconfig string) {

	dir := filepath.Join(composeConfig.ContentWorkingDir, monakoMenuDirectory)
	dst := filepath.Join(dir, "index.md")

	extractTheme(composeConfig)
	err := createHugoConfig(composeConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(dir, os.FileMode(0744))
	if err != nil {
		log.Fatalf("Error menu dir %s", err)
	}

	data, err := ioutil.ReadFile(menuconfig)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func createHugoConfig(composeConfig compose.Config) error {
	configContent := fmt.Sprintf(`# Autogenerated by Monako, do not edit
baseURL = '%s'
title = '%s'
theme = '%s'

# Because of this bug: https://github.com/gohugoio/hugo/issues/4841
# Maybe delete seems to be related to slow Github Actions
timeout = 60000

# Book configuration
disablePathToLower = true
enableGitInfo = false

# Needed for mermaid/katex shortcodes
[markup]
[markup.goldmark.renderer]
unsafe = true

[markup.tableOfContents]
startLevel = 1

[params]
# See: https://github.com/snipem/monako-book#configuration for settings
BookToC = true
BookLogo = '%s'
BookMenuBundle = '/%s'
BookSection = 'docs'
#BookRepo = 'https://github.com/alex-shpak/hugo-book'
#BookEditPath = 'edit/master/exampleSite/content'
BookDateFormat = 'Jan 2, 2006'
BookSearch = true
BookComments = true
	`, composeConfig.BaseURL, composeConfig.Title, themeName, composeConfig.Logo, monakoMenuDirectory)
	return ioutil.WriteFile(filepath.Join(composeConfig.HugoWorkingDir, "config.toml"), []byte(configContent), os.FileMode(0700))
}

func extractTheme(composeConfig compose.Config) {
	themezip, err := Asset("tmp/theme.zip")
	if err != nil {
		log.Fatalf("Error loading theme %s", err)
	}
	byteReader := bytes.NewReader(themezip)

	destPath, err := unpackit.Unpack(byteReader, filepath.Join(composeConfig.HugoWorkingDir, "themes"))
	if err != nil {
		log.Fatalf("Error extracting theme: %s", err)
	}

	log.Printf("Extracted %s", destPath)
}

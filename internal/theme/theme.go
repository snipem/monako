package theme

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/artdarek/go-unzip"
	"github.com/snipem/monako/pkg/compose"
)

const monakoMenuDirectory = "monako_menu_directory"
const themeName = "monako-book-master"

// CreateHugoPage extracts the Monako theme and copies the hugoconfig and menuconfig to the needed files
func CreateHugoPage(composeConfig compose.ComposeConfig, menuconfig string) {

	dir := filepath.Join(composeConfig.CompositionDir, composeConfig.ContentDir, monakoMenuDirectory)
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

func createHugoConfig(composeConfig compose.ComposeConfig) error {
	configContent := fmt.Sprintf(`# Autogenerated by Monako, do not edit
baseURL = '%s'
title = '%s'
theme = '%s'

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
BookLogo = 'logo.png'
BookMenuBundle = '/%s'
BookSection = 'docs'
#BookRepo = 'https://github.com/alex-shpak/hugo-book'
#BookEditPath = 'edit/master/exampleSite/content'
BookDateFormat = 'Jan 2, 2006'
BookSearch = true
BookComments = true
	`, composeConfig.BaseURL, composeConfig.Title, themeName, monakoMenuDirectory)
	return ioutil.WriteFile(filepath.Join(composeConfig.CompositionDir, "config.toml"), []byte(configContent), os.FileMode(0700))
}

func extractTheme(composeConfig compose.ComposeConfig) {
	themezip, err := Asset("tmp/theme.zip")
	if err != nil {
		log.Fatalf("Error loading theme %s", err)
	}

	// TODO Don't use local filesystem, keep it in memory
	tmpFile, err := ioutil.TempFile(os.TempDir(), "monako-theme-")
	if err != nil {
		log.Fatalf("Cannot create temporary file %s", err)
	}
	_, err = tmpFile.Write(themezip)
	if err != nil {
		log.Fatalf("Error temporary theme zip file %s", err)
	}

	tempfilename := tmpFile.Name()

	if err != nil {
		log.Fatalf("Error writing temp theme %s", err)
	}

	// TODO Don't use a library that depends on local files
	uz := unzip.New(tempfilename, filepath.Join(composeConfig.CompositionDir, "themes"))
	err = uz.Extract()
	if err != nil {
		log.Fatalf("Error extracting theme: %s ", err)
	}
	os.RemoveAll(tempfilename)
}

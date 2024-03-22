package compose

import (
	"fmt"
	theme "github.com/snipem/monako/pkg/compose/internal"
	"github.com/snipem/monako/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const monakoMenuDirectory = "monako_menu_directory"
const themeName = "monako-book"

// extractTheme extracts the Monako Theme to the Hugo Working Directory
func extractTheme(hugoWorkingDir string) error {
	themesDir := filepath.Join(hugoWorkingDir, "themes")
	// theme.RestoreAssets is autogeneration by the generation in cmd/monake/main.go
	// use "make theme" to generate those files
	err := theme.RestoreAssets(themesDir, themeName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error restoring asset %s to %s", themeName, themesDir))
	}
	return nil
}

// createMonakoStructureInHugoFolder extracts the Monako theme and copies the hugoconfig and menuconfig to the needed files
func createMonakoStructureInHugoFolder(composeConfig *Config, menuconfig string) error {

	var foldersToCreate = []string{"content", "themes"}
	for _, folder := range foldersToCreate {
		createDir := filepath.Join(composeConfig.ContentWorkingDir, folder)
		err := os.MkdirAll(createDir, standardFilemode)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error creating Monako structure %s", createDir))
		}
	}

	err := extractTheme(composeConfig.HugoWorkingDir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error extracting Hugo Theme"))
	}

	err = createHugoConfig(composeConfig)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating Hugo config"))
	}

	err = createMenuConfig(composeConfig, menuconfig)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating Monako menu config"))
	}
	return nil
}

func createMenuConfig(composeConfig *Config, menuconfig string) error {

	dir := filepath.Join(composeConfig.ContentWorkingDir, monakoMenuDirectory)
	dst := filepath.Join(dir, "index.md")
	err := os.MkdirAll(dir, os.FileMode(0744))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating menu dir %s", dir))
	}

	data, err := ioutil.ReadFile(menuconfig)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading menu config %s", menuconfig))
	}

	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error writing menu config to %s", dst))
	}
	return nil

}

// TODO Make MonakoGitLinks configurable

func createHugoConfig(composeConfig *Config) error {
	configContent := fmt.Sprintf(`# Autogenerated by Monako, do not edit
baseURL = '%s'
title = '%s'
theme = '%s'

# Use Uglyurls with html in path
uglyurls = true

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

[security]
  enableInlineShortcodes = false
  [security.exec]
    allow = ['^asciidoctor$', '^dart-sass-embedded$', '^go$', '^npx$', '^postcss$']
    osEnv = ['(?i)^(PATH|PATHEXT|APPDATA|TMP|TEMP|TERM)$']
  [security.funcs]
    getenv = ['^HUGO_']
  [security.http]
    methods = ['(?i)GET|POST']
    urls = ['.*']

[markup.asciidocext]
extensions = ["asciidoctor-diagram"]
workingFolderCurrent = true
# Use trace together with -v in hugo run
trace = false

[markup.asciidocext.attributes]
# this is needed for rendering section 0 to h1
showtitle = "true"

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

# Monako
MonakoGitLinks = true
MonakoDisableGitCommit = %v

	`, composeConfig.BaseURL, composeConfig.Title, themeName, composeConfig.Logo, monakoMenuDirectory, composeConfig.DisableCommitInfo)

	err := os.MkdirAll(composeConfig.HugoWorkingDir, standardFilemode)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(composeConfig.HugoWorkingDir, "config.toml"), []byte(configContent), standardFilemode)
	if err != nil {
		return err
	}

	return nil
}

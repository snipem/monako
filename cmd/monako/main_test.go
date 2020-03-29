package main

// run: go test -v ./cmd/monako

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	targetDir := filet.TmpDir(t, "")

	os.Args = []string{
		"monako",
		"-fail-on-error",
		"-target-dir", targetDir,
		"-config", "../../test/configs/only_markdown/config.markdown.yaml",
		"-menu-config", "../../test/configs/only_markdown/config.menu.markdown.md"}
	main()

	assert.FileExists(t, filepath.Join(targetDir, "compose/config.toml"), "Hugo config is not present")
	assert.FileExists(t, filepath.Join(targetDir, "compose/content/monako_menu_directory/index.md"), "Menu is not present")
	assert.FileExists(t, filepath.Join(targetDir, "compose/public/docs/test/test_doc_markdown/index.html"), "Generated Test doc does not exist")

	contentBytes, err := ioutil.ReadFile(filepath.Join(targetDir, "compose/public/docs/test/test_doc_markdown/index.html"))
	content := string(contentBytes)

	assert.NoError(t, err, "Can't read file")

	assert.Contains(t, content, "<strong>Test docs</strong>", "Contains menu")

	assert.Contains(t, content, "<img src=\"../profile.png\" alt=\"Picture in same folder\" />", "Contains relative picture")
	assert.FileExists(t, filepath.Join(targetDir, "compose/public/docs/test/profile.png"), "Relative picture right placed")

	assert.FileExists(t, filepath.Join(targetDir, "compose/public/docs/test/subfolder/subfolderprofile.png"), "Relative subfolder picture right placed")
	assert.Contains(t, content, "<img src=\"../subfolder/subfolderprofile.png\" alt=\"Picture in sub folder\" />", "Contains relative picture")

	if !t.Failed() {
		// Don't clean up when failed
		filet.CleanUp(t)
	}
}

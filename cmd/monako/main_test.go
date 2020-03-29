package main

// run: go test -v ./cmd/monako

import (
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

	if !t.Failed() {
		// Don't clean up when failed
		filet.CleanUp(t)
	}
}

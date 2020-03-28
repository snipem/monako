package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{
		"monako",
		"-fail-on-error",
		"-target-dir", filepath.Join(os.TempDir(), t.Name()),
		"-config", "../../test/configs/only_markdown/config.markdown.yaml",
		"-menu-config", "../../test/configs/only_markdown/config.menu.markdown.md"}
	main()
}

package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{
		"monako",
		"-fail-on-error",
		"-target-dir", "/tmp/target",
		"-config", "../../test/configs/only_markdown/config.markdown.yaml",
		"-menu-config", "../../test/configs/only_markdown/config.menu.markdown.md"}
	main()
}

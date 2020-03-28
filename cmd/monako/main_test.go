package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{
	"monako",
	// "-fail-on-error",
	"-target-dir", "/tmp/target",
    "-config", "../../test/config.local.yaml",
	"-menu-config", "../../test/config.menu.local.md"}
	main()	
}
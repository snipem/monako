package helpers

// run: make test

import (
	"os"

	hugo "github.com/gohugoio/hugo/commands"
)

// CleanUp removes the compose folder
func CleanUp() {
	os.RemoveAll("compose")
}

// HugoRun runs Hugo like the command line interface
func HugoRun(args []string) {
	hugo.Execute(args)
}

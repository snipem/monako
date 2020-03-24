package helpers

// run: make run

import (
	"os"

	hugo "github.com/gohugoio/hugo/commands"
)

// CleanUp removes the compose folder
func CleanUp() {
	os.RemoveAll("compose")
}

// HugoRun runs Hugo like the command line interface
func HugoRun(args []string) error {
	response := hugo.Execute(args)
	return response.Err
}

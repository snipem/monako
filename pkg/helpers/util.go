package helpers

// run: make run

import (
	hugo "github.com/gohugoio/hugo/commands"
)

// HugoRun runs Hugo like the command line interface
func HugoRun(args []string) error {
	response := hugo.Execute(args)
	return response.Err
}

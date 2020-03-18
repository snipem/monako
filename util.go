package main

// run: make test

import (
	"os"

	"github.com/gohugoio/hugo/commands"
)

func cleanUp() {
	os.RemoveAll("compose")
}

func hugoRun(args []string) {
	commands.Execute(args)
}

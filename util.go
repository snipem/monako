package main

// run: make test

import (
	"os"

	"github.com/gohugoio/hugo/commands"
)

func CleanUp() {
	os.RemoveAll("compose")
}

func HugoRun(args []string) {
	// args := []string{"--contentDir", "compose"}
	commands.Execute(args)
}

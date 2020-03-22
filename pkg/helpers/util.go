package helpers

// run: make test

import (
	"os"

	"github.com/gohugoio/hugo/commands"
)

func CleanUp() {
	os.RemoveAll("compose")
}

func HugoRun(args []string) {
	commands.Execute(args)
}

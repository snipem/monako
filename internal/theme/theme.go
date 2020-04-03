package theme

import (
	"bytes"
	"log"
	"path/filepath"

	"github.com/c4milo/unpackit"
)

func ExtractTheme(hugoWorkingDir string) {
	themezip, err := Asset("tmp/theme.zip")
	if err != nil {
		log.Fatalf("Error loading theme %s", err)
	}
	byteReader := bytes.NewReader(themezip)

	destPath, err := unpackit.Unpack(byteReader, filepath.Join(hugoWorkingDir, "themes"))
	if err != nil {
		log.Fatalf("Error extracting theme: %s", err)
	}

	log.Printf("Extracted %s", destPath)
}

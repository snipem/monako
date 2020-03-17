package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/artdarek/go-unzip"
	"github.com/codeskyblue/go-sh"
)

func getTheme(hugoconfig string, menuconfig string) {

	extractTheme()
	// TODO has to be TOML
	sh.Command("cp", hugoconfig, "compose/config.toml").Run()

	sh.Command("mkdir", "-p", "compose/content/menu/").Run()
	sh.Command("cp", menuconfig, "compose/content/menu/index.md").Run()

	addFixForADocTocToTheme()
}

func extractTheme() {
	themezip, err := Asset("tmp/theme.zip")
	if err != nil {
		log.Fatalf("Error loading theme %s", err)
	}

	// TODO Don't use local filesystem, keep it in memory
	tmpFile, err := ioutil.TempFile(os.TempDir(), "monako-theme-")
	if err != nil {
		fmt.Println("Cannot create temporary file", err)
	}
	tmpFile.Write(themezip)
	tempfilename := tmpFile.Name()

	if err != nil {
		log.Fatalf("Error writing temp theme %s", err)
	}

	// TODO Don't use a library that depends on local files
	uz := unzip.New(tempfilename, "compose/themes")
	err = uz.Extract()
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll(tempfilename)
}

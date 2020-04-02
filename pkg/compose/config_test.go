package compose

// run: MONAKO_TEST_REPO="/tmp/testdata/monako-test" go test ./pkg/compose/
// run: go test -v ./pkg/compose/

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
)

func TestGetTestConfig(t *testing.T) {
	config, workingdir, err := getTestConfig(t)
	assert.DirExists(t, workingdir)
	assert.NoError(t, err)
	assert.NotNil(t, config)

	t.Run("Check if init was right", func(t *testing.T) {

		assert.NotNil(t, config.HugoWorkingDir)
		assert.NotNil(t, config.ContentWorkingDir)

	})

	t.Run("Check if origins have pointer to config", func(t *testing.T) {

		for _, origin := range config.Origins {
			assert.NotNil(t, origin.config)
		}

	})
}

func TestCompose(t *testing.T) {
	config, _, err := getTestConfig(t)
	assert.NoError(t, err)

	config.Compose()

	wantFiles := []string{
		"docs/monako-test/README.md",
		"docs/monako-test/profile.png",
		"docs/monako-test/subfolder/subfolderprofile.png",
		"docs/monako-test/subfolder/test_doc_asciidoc_include_me_subfolder.adoc",
	}

	for _, wantFile := range wantFiles {
		assert.FileExists(t, filepath.Join(config.ContentWorkingDir, wantFile))
	}

}

func getTestConfig(t *testing.T) (config *Config, tempdir string, err error) {

	var testRepo string

	if os.Getenv("MONAKO_TEST_REPO") != "" {
		testRepo = os.Getenv("MONAKO_TEST_REPO")
	} else {
		testRepo = "https://github.com/snipem/monako-test.git"
	}

	tempdir = filet.TmpDir(t, os.TempDir())

	config = &Config{
		BaseURL:       "http://exampleurl.com",
		FileWhitelist: []string{".md", ".adoc", ".png"},
		Title:         "Test Config Title",
		Origins: []Origin{
			*NewOrigin(testRepo, "master", ".", "docs/monako-test"),
		},
	}

	config.initConfig(tempdir)

	return config, tempdir, nil
}

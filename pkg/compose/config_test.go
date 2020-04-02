package compose

// run: go test -v ./pkg/compose/

import (
	"os"
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

func getTestConfig(t *testing.T) (config *Config, tempdir string, err error) {

	tempdir = filet.TmpDir(t, os.TempDir())

	config = &Config{
		BaseURL:       "http://exampleurl.com",
		FileWhitelist: []string{".md", ".adoc", ".png"},
		Title:         "Test Config Title",
		Origins: []Origin{
			*NewOrigin("https://github.com/snipem/monako-test.git", "master", ".", "docs/monako-test"),
		},
	}

	config.initConfig(tempdir)

	return config, tempdir, nil
}

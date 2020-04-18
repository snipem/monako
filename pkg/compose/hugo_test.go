package compose

// run: go test  ./pkg/compose -run TestCreatePage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
)

func TestCreatePage(t *testing.T) {

	config, _ := getTestConfig(t)

	t.Run("Create Config", func(t *testing.T) {
		err := createHugoConfig(config)
		assert.NoError(t, err)
		t.Logf("Create in Hugo Working dir %s", config.HugoWorkingDir)

		// Check if Hugo content dir has been created
		assert.FileExists(t, filepath.Join(config.HugoWorkingDir, "config.toml"))
	})

	t.Run("Create Monako structure", func(t *testing.T) {
		testMenuConfig := filet.TmpFile(t, os.TempDir(), "# Test Config")
		err := createMonakoStructureInHugoFolder(config, testMenuConfig.Name())
		assert.NoError(t, err)

		// Check if Monako menu is in place
		assert.FileExists(t, filepath.Join(config.ContentWorkingDir, monakoMenuDirectory, "index.md"))

		// Check if theme has been extracted
		assert.FileExists(t, filepath.Join(config.HugoWorkingDir, "themes", themeName, "theme.toml"))

	})

}

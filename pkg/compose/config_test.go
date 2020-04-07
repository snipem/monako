package compose

// run: MONAKO_TEST_REPO="/tmp/testdata/monako-test" go test -v ./pkg/compose/

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
)

func TestGetTestConfig(t *testing.T) {
	config, workingdir := getTestConfig(t)
	assert.DirExists(t, workingdir)
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
	config, _ := getTestConfig(t)

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

// getLocalOrRemoteRepo returns a local or remote test remote to https://github.com/snipem/monako-test.git
// depending on if the MONAKO_TEST_REPO env variable is set or not
func getLocalOrRemoteRepo(t *testing.T) *Origin {

	var testRepo string

	if os.Getenv("MONAKO_TEST_REPO") != "" {
		testRepo = os.Getenv("MONAKO_TEST_REPO")
		t.Logf("Using local test repo: %s", testRepo)
	} else {
		testRepo = "https://github.com/snipem/monako-test.git"
	}
	return NewOrigin(testRepo, "master", ".", "docs/monako-test")

}

// getTestConfig returns a test config with a variable list of origins. If no origin is set
// as a parameter an example configuration is returned
func getTestConfig(t *testing.T, origins ...Origin) (config *Config, tempdir string) {

	var testOrigins []Origin

	if origins == nil {
		testOrigins = append(testOrigins, *getLocalOrRemoteRepo(t))
	} else {
		testOrigins = origins
	}

	tempdir = filet.TmpDir(t, os.TempDir())

	config = &Config{
		BaseURL:       "http://exampleurl.com",
		FileWhitelist: []string{".md", ".adoc", ".png"},
		Title:         "Test Config Title",
		Origins:       testOrigins,
	}

	config.initConfig(tempdir)

	assert.DirExists(t, tempdir)

	return config, tempdir
}

func TestInit(t *testing.T) {
	localFolder := "tmp/testdata"
	commandLineBaseURL := "http://overwrite.config"
	menuConfigFile := filet.TmpFile(t, os.TempDir(), "# Empty Menu")

	config := Init(CommandLineSettings{
		ConfigFilePath:     "../../test/config.local.yaml",
		MenuConfigFilePath: menuConfigFile.Name(),
		BaseURL:            commandLineBaseURL,
		ContentWorkingDir:  localFolder,
		FailOnHugoError:    true,
		Trace:              true,
	})

	assert.NotNil(t, config)
	assert.Equal(t, commandLineBaseURL, config.BaseURL)

	t.Run("Run Hugo", func(t *testing.T) {

		err := config.Run()
		assert.NoError(t, err)

	})

}

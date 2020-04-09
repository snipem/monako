package compose

// run: MONAKO_HUGE_REPOS_TEST=true go test -v ./pkg/compose/ -run TestHugeRepositories

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Flaque/filet"
	"github.com/snipem/monako/pkg/helpers"
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

func TestHugeRepositories(t *testing.T) {

	if os.Getenv("MONAKO_HUGE_REPOS_TEST") == "" {
		t.Skip("MONAKO_HUGE_REPOS_TEST is not set")
	}

	helpers.Trace()

	start := time.Now()
	config, _ := getTestConfig(t, *NewOrigin(
		// local path is $HOME/temp/hugo for hugo source code with lots of commits
		filepath.Join(os.Getenv("HOME"), "temp", "hugo"),
		"master",
		"",
		"huge/test/docs",
	))

	config.Compose()

	t.Logf("took %v\n", time.Since(start))

}

func TestDeactivatedCommitInfo(t *testing.T) {
	config, _ := getTestConfig(t)
	// Standard must be false
	assert.False(t, config.DisableCommitInfo)
	config.DisableCommitInfo = true

	config.Compose()
	for i := range config.Origins[0].Files {
		// Commit must be nil
		assert.Nil(t, config.Origins[0].Files[i].Commit)
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

	t.Run("Generate HTML with Hugo", func(t *testing.T) {

		err := config.Generate()
		assert.NoError(t, err)

	})

}

package compose

// run: make benchmark

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

	err := config.Compose()
	assert.NoError(t, err)

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

func BenchmarkHugoRepositorySingleFiles(b *testing.B) {

	origin := NewOrigin(
		// local path is $HOME/temp/hugo for hugo source code with lots of commits
		filepath.Join(os.Getenv("HOME"), "/temp/monako-testrepos/hugo"),
		"master",
		"",
		"huge/test/docs",
	)

	origin.config = &Config{
		DisableCommitInfo: false,
	}

	_, err := origin.CloneDir()
	assert.NoError(b, err)

	// Don't fetch commit info early
	origin.config = &Config{DisableCommitInfo: true}
	b.Run("Get Commit Info for Hugo", func(b *testing.B) {

		for n := 0; n < b.N; n++ {
			_, err := getCommitInfo("README.md", origin.repo)
			assert.NoError(b, err)

			// Older commit long time no change, far behind in git log
			_, err = getCommitInfo("docs/archetypes/default.md", origin.repo)
			assert.NoError(b, err)
		}

	})

}

func BenchmarkSlowRepositorySingleFiles(b *testing.B) {

	if _, isSet := os.LookupEnv("MONAKO_SLOW_REPO_FOLDER"); !isSet {
		assert.FailNow(b, "Env var MONAKO_SLOW_REPO_FOLDER not set")
	}

	if _, isSet := os.LookupEnv("MONAKO_SLOW_REPO_FILE1"); !isSet {
		assert.FailNow(b, "Env var MONAKO_SLOW_REPO_FILE1 not set")
	}

	if _, isSet := os.LookupEnv("MONAKO_SLOW_REPO_FILE2"); !isSet {
		assert.FailNow(b, "Env var MONAKO_SLOW_REPO_FILE2 not set")
	}

	// Use env vars, these folders are secret because of client project
	slowRepoFolder := os.Getenv("MONAKO_SLOW_REPO_FOLDER")
	slowRepoFile1 := os.Getenv("MONAKO_SLOW_REPO_FILE1")
	slowRepoFile2 := os.Getenv("MONAKO_SLOW_REPO_FILE2")

	origin := NewOrigin(
		// local path is $HOME/temp/hugo for hugo source code with lots of commits
		filepath.Join(os.Getenv("HOME"), slowRepoFolder),
		"develop",
		"",
		"../../tmp/testdata/huge/test/docs",
	)

	// Don't fetch commit info early
	origin.config = &Config{DisableCommitInfo: true}
	_, err := origin.CloneDir()
	assert.NoError(b, err)

	b.Run("Get Commit Info", func(b *testing.B) {

		for n := 0; n < b.N; n++ {

			_, err := getCommitInfo(slowRepoFile1, origin.repo)
			assert.NoError(b, err)

			// Older commit long time no change, far behind in git log
			_, err = getCommitInfo(slowRepoFile2, origin.repo)
			assert.NoError(b, err)
		}

	})

}

func BenchmarkWholeRepoHugoRepositoryWholeRepo(b *testing.B) {

	origin := NewOrigin(
		// local path is $HOME/temp/hugo for hugo source code with lots of commits
		filepath.Join(os.Getenv("HOME"), "/temp/monako-testrepos/hugo"),
		"master",
		"docs/content/en/commands",
		"../../tmp/testdata/huge/test/docs/hugo/",
	)

	origin.config = &Config{
		DisableCommitInfo: false,
		FileWhitelist:     []string{".md", ".png"},
	}

	origin.FileWhitelist = origin.config.FileWhitelist
	filesystem, err := origin.CloneDir()
	assert.NoError(b, err)

	// Don't fetch commit info early
	b.Run("Get Commit Info for Hugo", func(b *testing.B) {

		for n := 0; n < b.N; n++ {
			err := origin.ComposeDir(filesystem)
			assert.NoError(b, err)
		}

	})

}

func BenchmarkWholeRepoSlowRepository(b *testing.B) {

	if _, isSet := os.LookupEnv("MONAKO_SLOW_REPO_FOLDER"); !isSet {
		assert.FailNow(b, "Env var MONAKO_SLOW_REPO_FOLDER not set")
	}

	if _, isSet := os.LookupEnv("MONAKO_SLOW_REPO_SOURCE"); !isSet {
		assert.FailNow(b, "Env var MONAKO_SLOW_REPO_SOURCE not set")
	}

	// Use env vars, these folders are secret because of a client project
	slowRepoFolder := os.Getenv("MONAKO_SLOW_REPO_FOLDER")
	slowRepoSource := os.Getenv("MONAKO_SLOW_REPO_SOURCE")

	origin := NewOrigin(
		// local path is $HOME/temp/hugo for hugo source code with lots of commits
		filepath.Join(os.Getenv("HOME"), slowRepoFolder),
		"develop",
		slowRepoSource,
		"../../tmp/testdata/huge/test/docs/slow",
	)

	origin.config = &Config{
		DisableCommitInfo: false,
		FileWhitelist:     []string{".md", ".png"},
	}

	origin.FileWhitelist = origin.config.FileWhitelist

	filesystem, err := origin.CloneDir()
	assert.NoError(b, err)
	b.Run("Get Commit Info", func(b *testing.B) {

		for n := 0; n < b.N; n++ {
			err := origin.ComposeDir(filesystem)
			assert.NoError(b, err)
		}

	})

}

func TestDeactivatedCommitInfo(t *testing.T) {
	config, _ := getTestConfig(t)
	// Standard must be false
	assert.False(t, config.DisableCommitInfo)
	config.DisableCommitInfo = true

	err := config.Compose()
	assert.NoError(t, err)

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
func GetLocalTempDir(t *testing.T) (tempdir string) {

	localTmpDir := filepath.Join("../../tmp/testdata/", t.Name())
	err := os.MkdirAll(localTmpDir, standardFilemode)
	assert.NoError(t, err)

	return filet.TmpDir(t, localTmpDir)
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

	tempdir = GetLocalTempDir(t)

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
	localFolder := GetLocalTempDir(t)
	commandLineBaseURL := "http://overwrite.config"
	menuConfigFile := filet.TmpFile(t, os.TempDir(), "# Empty Menu")

	config := Init(CommandLineSettings{
		ConfigFilePath:     "../../test/config.local.yaml",
		MenuConfigFilePath: menuConfigFile.Name(),
		BaseURL:            commandLineBaseURL,
		ContentWorkingDir:  localFolder,
		FailOnHugoError:    true,
	})

	assert.NotNil(t, config)
	assert.Equal(t, commandLineBaseURL, config.BaseURL)

	t.Run("Generate HTML with Hugo", func(t *testing.T) {

		err := config.Generate()
		assert.NoError(t, err)

	})

}

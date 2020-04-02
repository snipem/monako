package compose

// run: go test -v ./pkg/compose

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Load good Config", func(t *testing.T) {
		customWorkingdir := "workingdir123"
		config, err := LoadConfig("../../test/config.local.yaml", customWorkingdir)
		assert.NoError(t, err)

		assert.Contains(t, config.ContentWorkingDir, customWorkingdir)
		assert.NotEmpty(t, config.HugoWorkingDir, customWorkingdir)

		assert.NotEmpty(t, config.Origins)

	})

	t.Run("Fail on missing config", func(t *testing.T) {
		_, err := LoadConfig("missing path", "")
		assert.Error(t, err)
	})

}

func TestCleanUp(t *testing.T) {
	config, _, err := getTestConfig(t)
	assert.NoError(t, err)

	tmpFile := filepath.Join(
		config.HugoWorkingDir,
		"testfile.txt")

	// Create the test data because it is not existing yet
	err = os.Mkdir(config.HugoWorkingDir, standardFilemode)
	if err != nil {
		log.Fatal(err)
	}

	t.Logf("Using temp file %s", tmpFile)
	_ = ioutil.WriteFile(tmpFile, []byte("none"), standardFilemode)

	assert.FileExists(t, tmpFile, "File is existing that is to be cleaned up")
	config.CleanUp()
	assert.NoFileExists(t, tmpFile, "File seems not to be cleaned up, is stil present")

}

// TestLocalPath tests if the local file path calculation for remote files is correct
func TestLocalPath(t *testing.T) {

	equalPath(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", ".", ".", "filename.md"),
		"Simple setup, always first level")

	equalPath(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", "docs", ".", "docs/filename.md"),
		"With remote 'docs' folder")

	equalPath(t,
		"/tmp/compose/docs/filename.md",
		getLocalFilePath("/tmp/compose", ".", ".", "docs/filename.md"),
		"With remote 'docs' folder, but keep structure")

	equalPath(t,
		"compose/filename.md",
		getLocalFilePath("./compose", ".", ".", "filename.md"),
		"Path is relative")

	equalPath(t,
		"/tmp/compose/localTarget/filename.md",
		getLocalFilePath("/tmp/compose", ".", "localTarget", "filename.md"),
		"Local Target folder")

	equalPath(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", ".", "", "filename.md"),
		"Empty local target folder")
}

// equalPath is like assert.Equal but with ignoring operation system specifc pathes.
// On Unix "/" and Windows "\" systems this check compares pathes either way.
func equalPath(t *testing.T, expected string, actual string, msg string) {

	assert.Equal(t,
		filepath.ToSlash(expected),
		filepath.ToSlash(actual),
		msg,
	)
}

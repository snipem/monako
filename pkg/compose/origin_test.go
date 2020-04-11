package compose

// run: MONAKO_TEST_REPO="/tmp/testdata/monako-test" go test ./pkg/compose/

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugofs/files"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Load good Config", func(t *testing.T) {

		customWorkingdir := GetLocalTempDir(t)

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
	config, _ := getTestConfig(t)

	tmpFile := filepath.Join(
		config.HugoWorkingDir,
		"testfile.txt")

	// Create the test data because it is not existing yet
	err := os.Mkdir(config.HugoWorkingDir, standardFilemode)
	if err != nil {
		log.Fatal(err)
	}

	t.Logf("Using temp file %s", tmpFile)
	_ = ioutil.WriteFile(tmpFile, []byte("none"), standardFilemode)

	assert.FileExists(t, tmpFile, "File is existing that is to be cleaned up")
	config.CleanUp()
	assert.NoFileExists(t, tmpFile, "File seems not to be cleaned up, is stil present")

}

func TestGitCommiter(t *testing.T) {

	config, _ := getTestConfig(t)
	config.Compose()
	origins := config.Origins
	firstOrigin := origins[0]

	for i := 0; i < len(firstOrigin.Files); i++ {
		t.Run("Retrieve commit info of file", func(t *testing.T) {
			if files.IsContentFile(firstOrigin.Files[i].RemotePath) {
				ci := firstOrigin.Files[0].Commit
				assert.Contains(t, ci.Committer.Email, "@")
			} else {
				t.Logf("Skipping commit check for %s, is not a content file", firstOrigin.Files[i].RemotePath)
			}

		})
	}
}

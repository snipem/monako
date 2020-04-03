package compose

// run: MONAKO_TEST_REPO="/tmp/testdata/monako-test" go test ./pkg/compose/

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// var o *Origin

// func TestMain(m *testing.M) {
// 	// Setup git clone of repo
// 	setup()
// 	os.Exit(m.Run())
// }

// func setup() {
// 	o = NewOrigin("https://github.com/snipem/monako-test.git", "master", ".", "test/docs")
// 	o.CloneDir()
// }

// TestCopyDir is a test for testing the copying capability of a single directory
// func TestCopyDir(t *testing.T) {
// 	t.Skip("Skip for now, is hard to understand")

// 	// TODO Get own temporary test folder
// 	// targetDir := filepath.Join(os.TempDir(), "tmp/testrun/", t.Name())
// 	// defer os.RemoveAll(targetDir)

// 	var whitelist = []string{".md", ".png"}

// 	t.Run("start in single directory 'test'", func(t *testing.T) {
// 		o.SourceDir = "test"
// 		o.FileWhitelist = whitelist
// 		tempDir := filet.TmpDir(t, os.TempDir())
// 		o.ComposeDir()
// 		expectedTargetFile := filepath.Join(tempDir, "compose", "test_docs/test_doc_markdown.md")
// 		b, err := ioutil.ReadFile(expectedTargetFile)

// 		assert.NoError(t, err, "File not found")
// 		assert.Contains(t, string(b), "# Markdown Doc 1")
// 	})

// 	t.Run("start in deeper directory 'test/test_docs/'", func(t *testing.T) {
// 		o.SourceDir = "test/test_docs/"
// 		o.FileWhitelist = whitelist
// 		tempDir := filet.TmpDir(t, os.TempDir())
// 		o.ComposeDir()
// 		expectedTargetFile := filepath.Join(tempDir, "compose", "/test_doc_markdown.md")
// 		b, err := ioutil.ReadFile(expectedTargetFile)

// 		assert.NoError(t, err, "File not found")
// 		assert.Contains(t, string(b), "# Markdown Doc 1")
// 	})

// }

func TestGitCommiter(t *testing.T) {

	config, _, err := getTestConfig(t)
	assert.NoError(t, err)
	config.Compose()
	origins := config.Origins
	firstOrigin := origins[0]

	t.Run("Retrieve info of first file", func(t *testing.T) {

		assert.NotNil(t, firstOrigin.Files)
		ci := firstOrigin.Files[0].Commit

		assert.NoError(t, err, "Could not retrieve commit info")
		assert.Contains(t, ci.Committer.Email, "@")

	})

	t.Run("Second file", func(t *testing.T) {
		ci := firstOrigin.Files[1].Commit

		assert.NoError(t, err, "Could not retrieve commit info")
		assert.Contains(t, ci.Committer.Email, "@")

	})

	t.Run("Not existing file", func(t *testing.T) {

		t.Skip("This wont work right now")

	})

}

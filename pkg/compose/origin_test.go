package compose

// run: MONAKO_TEST_REPO="/tmp/testdata/monako-test" go test ./pkg/compose/

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCommiter(t *testing.T) {

	config, _ := getTestConfig(t)
	config.Compose()
	origins := config.Origins
	firstOrigin := origins[0]

	t.Run("Retrieve info of first file", func(t *testing.T) {

		assert.NotNil(t, firstOrigin.Files)
		ci := firstOrigin.Files[0].Commit

		assert.Contains(t, ci.Committer.Email, "@")

	})

	t.Run("Second file", func(t *testing.T) {
		ci := firstOrigin.Files[1].Commit

		assert.Contains(t, ci.Committer.Email, "@")

	})

	t.Run("Not existing file", func(t *testing.T) {

		t.Skip("This wont work right now")

	})

}

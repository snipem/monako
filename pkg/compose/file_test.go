package compose

// run: go test -v ./pkg/compose

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestLocalPath(t *testing.T) {

	assert.Equal(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", ".", ".", "filename.md"),
		"Simple setup, always first level")

	assert.Equal(t,
		"/tmp/compose/filename.md",
		getLocalFilePath("/tmp/compose", "docs", ".", "docs/filename.md"),
		"With remote 'docs' folder")

	assert.Equal(t,
		"/tmp/compose/docs/filename.md",
		getLocalFilePath("/tmp/compose", ".", ".", "docs/filename.md"),
		"With remote 'docs' folder, but keep structure")

	assert.Equal(t,
		"compose/filename.md",
		getLocalFilePath("./compose", ".", ".", "filename.md"),
		"Path is relative")

	assert.Equal(t,
		"/tmp/compose/localTarget/filename.md",
		getLocalFilePath("/tmp/compose", ".", "localTarget", "filename.md"),
		"Simple setup, always first level")
}

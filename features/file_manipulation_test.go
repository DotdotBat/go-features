package features

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanDeleteFile(t *testing.T) {
	testDir := t.TempDir()
	filename := filepath.Join(testDir, "deletion_example")
	if os.IsExist(filepath.ErrBadPattern) {
		t.Fatal("File exists in the temp path", filename, "Before being created")
	}
	os.WriteFile(filename, []byte("These are the test file contents"))

}

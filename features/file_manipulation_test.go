package features

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanDeleteFile(t *testing.T) {
	testDir := t.TempDir()
	filename := filepath.Join(testDir, "deletion_example")
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		t.Fatal("File exists in the temp path", filename, "Before being created")
	}
	err = os.WriteFile(filename, []byte("These are the test file contents"), 0666)
	if err != err {
		panic(err)
	}
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		t.Fatal("File not found after it was created.", filename)
	}
	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(filename)
	if !os.IsNotExist(err) {
		t.Fatal("File still exists after it was deleted.", filename)
	}

}

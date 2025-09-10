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

func TestCanMoveFile(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()
	filename := "move"
	sourceFilename := filepath.Join(sourceDir, filename)
	targetFilename := filepath.Join(targetDir, filename)
	err := os.WriteFile(sourceFilename, []byte("a file to be moved"), 0666)
	if err != nil {
		panic(err)
	}
	err = os.Rename(sourceFilename, targetFilename)
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(sourceFilename)
	if !os.IsNotExist(err) {
		t.Fatal("file still exists in source directory after move", sourceFilename)
	}
	_, err = os.Stat(targetFilename)
	if os.IsNotExist(err) {
		t.Fatal("File doesn't exists in the target directory after move")
	}
}

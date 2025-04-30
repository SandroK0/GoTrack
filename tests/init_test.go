package tests

import (
	"os"
	"path/filepath"
	"testing"

	"GoTrack/constants"
	"GoTrack/vcs"
)

func TestHandleInit_CreatesDirectories(t *testing.T) {
	tmp := t.TempDir()

	vcs.HandleInit(tmp)

	gtPath := filepath.Join(tmp, constants.GTDir)
	objectsPath := filepath.Join(tmp, constants.ObjectsDir)

	if _, err := os.Stat(gtPath); os.IsNotExist(err) {
		t.Fatalf("Expected .gt directory at %s, but it was not created", gtPath)
	}

	if _, err := os.Stat(objectsPath); os.IsNotExist(err) {
		t.Fatalf("Expected objects directory at %s, but it was not created", objectsPath)
	}

}

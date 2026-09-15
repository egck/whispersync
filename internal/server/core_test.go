package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/egck/whispersync/internal/config"
	"github.com/egck/whispersync/internal/types"
)

func TestCoreSyncFileCreatesFileAndPreviousVersion(t *testing.T) {
	storagePath := t.TempDir()
	core := &Core{
		AppConfig: &config.ServerConfig{DataStoragePath: storagePath},
	}

	firstRequest := &types.FileSyncRequest{
		SyncRequest: types.SyncRequest{FilePath: "/source/document.txt"},
		FileContent: []byte("first version"),
	}
	if err := core.SyncFile(firstRequest); err != nil {
		t.Fatalf("synchronize first version: %v", err)
	}

	currentPath := filepath.Join(storagePath, "document.txt")
	assertFileContent(t, currentPath, "first version")

	secondRequest := &types.FileSyncRequest{
		SyncRequest: types.SyncRequest{FilePath: "/source/document.txt"},
		FileContent: []byte("second version"),
	}
	if err := core.SyncFile(secondRequest); err != nil {
		t.Fatalf("synchronize second version: %v", err)
	}

	assertFileContent(t, currentPath, "second version")

	versionEntries, err := os.ReadDir(filepath.Join(storagePath, ".versions"))
	if err != nil {
		t.Fatalf("read versions directory: %v", err)
	}
	if len(versionEntries) != 1 {
		t.Fatalf("expected one stored version, got %d", len(versionEntries))
	}

	versionPath := filepath.Join(storagePath, ".versions", versionEntries[0].Name())
	assertFileContent(t, versionPath, "first version")
}

func assertFileContent(t *testing.T, filePath string, expectedContent string) {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read file %s: %v", filePath, err)
	}
	if string(content) != expectedContent {
		t.Fatalf("expected content %q, got %q", expectedContent, content)
	}
}

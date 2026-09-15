package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egck/whispersync/internal/types"
)

func TestRequestManagerSync(t *testing.T) {
	expectedContent := []byte("content to synchronize")
	receivedRequest := make(chan types.FileSyncRequest, 1)

	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", request.Method)
		}
		if request.URL.Path != "/api/sync/file" {
			t.Errorf("expected /api/sync/file path, got %s", request.URL.Path)
		}

		var syncRequest types.FileSyncRequest
		if err := json.NewDecoder(request.Body).Decode(&syncRequest); err != nil {
			t.Errorf("decode sync request: %v", err)
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		receivedRequest <- syncRequest
		responseWriter.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	filePath := filepath.Join(t.TempDir(), "document.txt")
	if err := os.WriteFile(filePath, expectedContent, 0600); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	requestManager := NewRequestManager(strings.TrimPrefix(server.URL, "http://"), "")
	if err := requestManager.Sync(filePath); err != nil {
		t.Fatalf("synchronize file: %v", err)
	}

	syncRequest := <-receivedRequest
	if syncRequest.FilePath != filePath {
		t.Fatalf("expected path %q, got %q", filePath, syncRequest.FilePath)
	}
	if string(syncRequest.FileContent) != string(expectedContent) {
		t.Fatalf("expected content %q, got %q", expectedContent, syncRequest.FileContent)
	}
}

func TestRequestManagerSyncReturnsServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, "storage unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	filePath := filepath.Join(t.TempDir(), "document.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0600); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	requestManager := NewRequestManager(strings.TrimPrefix(server.URL, "http://"), "")
	err := requestManager.Sync(filePath)
	if err == nil {
		t.Fatal("expected synchronization error, got nil")
	}
	if !strings.Contains(err.Error(), "storage unavailable") {
		t.Fatalf("expected server error body, got %q", err)
	}
}

func TestNewRequestManagerConfiguresHTTPTimeout(t *testing.T) {
	requestManager := NewRequestManager("localhost:8080", "")

	if requestManager.Client == nil {
		t.Fatal("expected HTTP client, got nil")
	}
	if requestManager.Client.Timeout != defaultHTTPClientTimeout {
		t.Fatalf("expected timeout %s, got %s", defaultHTTPClientTimeout, requestManager.Client.Timeout)
	}
}

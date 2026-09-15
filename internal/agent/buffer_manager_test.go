package agent

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/egck/whispersync/internal/types"
	"github.com/fsnotify/fsnotify"
)

func TestBufferManagerFlushDeduplicatesFileEvents(t *testing.T) {
	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCount.Add(1)
		responseWriter.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	filePath := createBufferedTestFile(t)
	requestManager := NewRequestManager(strings.TrimPrefix(server.URL, "http://"), "")
	bufferManager := NewBufferManager(1, types.NewAppContext(), requestManager)

	bufferManager.Write(&fsnotify.Event{Name: filePath, Op: fsnotify.Write})
	bufferManager.Write(&fsnotify.Event{Name: filePath, Op: fsnotify.Write})
	bufferManager.Flush()

	if requestCount.Load() != 1 {
		t.Fatalf("expected one synchronization request, got %d", requestCount.Load())
	}
}

func TestBufferManagerFlushRequeuesFailedEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	filePath := createBufferedTestFile(t)
	requestManager := NewRequestManager(strings.TrimPrefix(server.URL, "http://"), "")
	bufferManager := NewBufferManager(1, types.NewAppContext(), requestManager)
	originalEvent := &fsnotify.Event{Name: filePath, Op: fsnotify.Write}

	bufferManager.Write(originalEvent)
	bufferManager.Flush()

	bufferManager.lock.Lock()
	defer bufferManager.lock.Unlock()

	singleBuffer, exists := bufferManager.SingleBuffers[filePath]
	if !exists {
		t.Fatal("expected failed event to remain buffered")
	}
	if singleBuffer.event != originalEvent {
		t.Fatal("expected failed event to be queued for retry")
	}
}

func createBufferedTestFile(t *testing.T) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), "document.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0600); err != nil {
		t.Fatalf("create buffered test file: %v", err)
	}

	return filePath
}

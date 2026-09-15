package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egck/whispersync/internal/config"
	"github.com/egck/whispersync/internal/types"
	"github.com/gin-gonic/gin"
)

func TestSyncFileRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storagePath := t.TempDir()
	core := &Core{AppConfig: &config.ServerConfig{DataStoragePath: storagePath}}
	router := gin.New()
	registerRoutes(router, core)

	syncRequest := types.FileSyncRequest{
		SyncRequest: types.SyncRequest{FilePath: "/source/document.txt"},
		FileContent: []byte("synchronized content"),
	}
	requestBody, err := json.Marshal(syncRequest)
	if err != nil {
		t.Fatalf("marshal sync request: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/sync/file", bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	content, err := os.ReadFile(filepath.Join(storagePath, "document.txt"))
	if err != nil {
		t.Fatalf("read synchronized file: %v", err)
	}
	if string(content) != "synchronized content" {
		t.Fatalf("unexpected synchronized content %q", content)
	}
}

func TestSyncFileRouteRejectsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	core := &Core{AppConfig: &config.ServerConfig{DataStoragePath: t.TempDir()}}
	router := gin.New()
	registerRoutes(router, core)

	request := httptest.NewRequest(http.MethodPost, "/api/sync/file", bytes.NewBufferString("{"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestNewAPIConfiguresServerTimeouts(t *testing.T) {
	core := &Core{
		AppConfig:  &config.ServerConfig{APIEndpoint: ":0"},
		AppContext: types.NewAppContext(),
	}
	api := NewAPI(core)

	if api.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("expected read header timeout 5s, got %s", api.ReadHeaderTimeout)
	}
	if api.ReadTimeout != 30*time.Second {
		t.Fatalf("expected read timeout 30s, got %s", api.ReadTimeout)
	}
	if api.WriteTimeout != 30*time.Second {
		t.Fatalf("expected write timeout 30s, got %s", api.WriteTimeout)
	}
	if api.IdleTimeout != 60*time.Second {
		t.Fatalf("expected idle timeout 60s, got %s", api.IdleTimeout)
	}
}

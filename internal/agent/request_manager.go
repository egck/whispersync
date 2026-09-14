package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/egck/whispersync/internal/types"
)

type RequestManager struct {
	secret    string
	ServerURI string
	Client    *http.Client
}

const defaultHTTPClientTimeout = 30 * time.Second

func NewRequestManager(serverUri string, secret string) *RequestManager {
	return &RequestManager{
		secret:    secret,
		ServerURI: serverUri,
		Client: &http.Client{
			Timeout: defaultHTTPClientTimeout,
		},
	}
}

func (requestManager *RequestManager) Sync(filePath string) error {
	// Prepare file for sync
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read file %s: %w", filePath, err)
	}
	syncRequest := types.FileSyncRequest{
		SyncRequest: types.SyncRequest{
			FilePath: filePath,
		},
		FileContent: fileContent,
	}
	jsonData, err := json.Marshal(syncRequest)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	// Post request
	response, err := requestManager.Client.Post("http://"+requestManager.ServerURI+"/api/sync/file", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("unable to post http request: %w", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("unable to close HTTP response body: %v", err)
		}
	}()
	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to read HTTP response body: %w", err)
		}
		return fmt.Errorf("unable to sync with server: %s", string(body))
	}

	return nil
}

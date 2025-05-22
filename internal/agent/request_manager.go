package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/egck/whispersync/internal/types"
)

type RequestManager struct {
	secret    string
	ServerURI string
}

func NewRequestManager(serverUri string, secret string) *RequestManager {
	return &RequestManager{
		secret:    secret,
		ServerURI: serverUri,
	}
}

func (requestManager *RequestManager) Sync(filePath string) {
	// Prepare file for sync
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("unable to sync: ", err.Error())
		return
	}
	syncRequest := types.FileSyncRequest{
		SyncRequest: types.SyncRequest{
			FilePath:       filePath,
		},
		FileContent: fileContent,
	}
	jsonData, err := json.Marshal(syncRequest)
	if err != nil {
		fmt.Println("unable to sync: ", err.Error())
		return
	}

	// Post request
	response, err := http.Post("http://"+requestManager.ServerURI+"/api/sync/file", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("unable to sync: ", err.Error())
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("unable to sync: unknown error")
		}
		fmt.Println("unable to sync with server: ", string(body))
	}
}

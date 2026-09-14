package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bxgstudio/goconfigloader/loader"
	"github.com/egck/whispersync/internal/config"
	"github.com/egck/whispersync/internal/types"
)

type Core struct {
	AppContext *types.AppContext
	AppConfig  *config.ServerConfig
	API        *API
}

func NewCore() *Core {
	core := &Core{
		AppContext: types.NewAppContext(),
		AppConfig:  &config.ServerConfig{},
	}

	// Load appConfig
	appConfigPath := os.Getenv("WS_SERVER_APP_CONFIG_PATH")
	if appConfigPath == "" {
		appConfigPath = "./app_config.yaml"
	}
	if err := loader.LoadConfig(appConfigPath, core.AppConfig); err != nil {
		log.Fatalf("unable to load appconfig file: %s", err)
	}

	// validate config
	if err := core.AppConfig.Validate(); err != nil {
		log.Fatalf("invalid appconfig: %s", err)
	}

	// Create backup path
	core.AppConfig.DataStoragePath = strings.TrimSuffix(core.AppConfig.DataStoragePath, "/")
	if err := os.MkdirAll(core.AppConfig.DataStoragePath, 0755); err != nil {
		fmt.Printf("unable to create backup dir at %s: %s\n", core.AppConfig.DataStoragePath, err.Error())
		os.Exit(1)
	}

	// Init api
	api := NewAPI(core)
	core.API = api

	return core
}

func (core *Core) Run() {
	core.AppContext.WaitGroup.Add(1)
	go func() {
		defer core.AppContext.WaitGroup.Done()
		core.API.Run()
	}()
}

func (core *Core) Stop() {
	fmt.Println("stopping app...")
	// Close core context
	core.AppContext.Stop()

	// Stop http server
	core.API.Stop()
	fmt.Println("application stopped properly")
}

func (core *Core) GetFileNames() []string {
	// Init file names slice
	fileNames := []string{}

	// Load all entries of the storage directory
	entries, err := os.ReadDir(core.AppConfig.DataStoragePath)
	if err != nil {
		fmt.Println("error while opening storage directory :", err)
		return fileNames
	}

	// Add only file names to slice
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	return fileNames
}

func (core *Core) SyncFile(syncRequest *types.FileSyncRequest) error {
	// Create files paths
	fileName := filepath.Base(syncRequest.FilePath)
	versionsPath := core.AppConfig.DataStoragePath + "/.versions"

	// Create version directory if needed
	if err := os.MkdirAll(versionsPath, 0755); err != nil {
		return err
	}

	// Write file and store old version if possible
	err := core.writeAndBackupVersion(core.AppConfig.DataStoragePath, versionsPath, fileName, syncRequest.FileContent)
	if err != nil {
		return err
	}
	return nil
}

func (core *Core) SyncRename(syncRequest *types.SyncRequest) error {
	fmt.Println(syncRequest)
	return nil
}

func (core *Core) writeAndBackupVersion(backupFilePath string, versionsPath string, fileName string, fileContent []byte) error {
	// Check if file already exists
	filePath := backupFilePath + "/" + fileName
	currentFile, err := os.OpenFile(filePath, os.O_RDWR, 0755)

	// Create new file if no anterior version exists or open existing file
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(filePath, fileContent, 0755); err != nil {
				return fmt.Errorf("unable to create file %s: %w", backupFilePath, err)
			}
			return nil
		}
		return fmt.Errorf("unable to open file %s: %w", backupFilePath, err)
	}
	defer func() {
		err := currentFile.Close()
		if err != nil {
			log.Printf("unable to close file %s: %v", filePath, err)
		}
	}()

	// Create a timestamped version of current file
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	versionedPath := versionsPath + "/" + fileName + "." + nowStr
	versionedFile, err := os.Create(versionedPath)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", versionedPath, err)
	}
	defer func() {
		err := versionedFile.Close()
		if err != nil {
			log.Printf("unable to close file %s: %v", versionedPath, err)
		}
	}()
	_, err = io.Copy(versionedFile, currentFile)
	if err != nil {
		return fmt.Errorf("unable to copy %s content in %s: %w", filePath, versionedPath, err)
	}

	// Update current file
	if err = currentFile.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate file %s: %w", filePath, err)
	}
	if _, err = currentFile.Seek(0, 0); err != nil {
		return fmt.Errorf("unable to reset %s file index: %w", filePath, err)
	}
	if _, err = currentFile.Write(fileContent); err != nil {
		return fmt.Errorf("unable to save file %s: %w", filePath, err)
	}

	// Flush files content to disk
	if err = currentFile.Sync(); err != nil {
		return fmt.Errorf("unable to flush %s on disk: %w", filePath, err)
	}
	if err = versionedFile.Sync(); err != nil {
		return fmt.Errorf("unable to flush %s on disk: %w", versionedPath, err)
	}

	return nil
}

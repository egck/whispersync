package agent

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/egck/whispersync/internal/types"
	"github.com/egck/whispersync/internal/utils"

	"github.com/fsnotify/fsnotify"
)

type WSWatcher struct {
	*fsnotify.Watcher
	AppContext    *types.AppContext
	bufferManager *BufferManager
}

func NewWatcher(AppContext *types.AppContext, directoryPath string, bufferManager *BufferManager) (*WSWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Add subdirectories
	err = filepath.WalkDir(directoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Println("watching: ", path)
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &WSWatcher{
		Watcher:       watcher,
		AppContext:    AppContext,
		bufferManager: bufferManager,
	}, nil
}

// TODO: Rework
func (watcher *WSWatcher) Run() {
	fmt.Println("starting watcher")
	defer watcher.Close()

	for {
		select {
		case <-watcher.AppContext.Context.Done():
			fmt.Println("watcher stopped properly")
			return
		case event := <-watcher.Events:
			// Process event
			watcher.ProcessEvent(&event)

		case err := <-watcher.Errors:
			log.Println("watcher error:", err)
		}
	}
}

func (watcher *WSWatcher) ignoreEvent(event *fsnotify.Event) bool {
	// ignore event that is not about creation, rename or write
	if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Rename) == 0 {
		return true
	}

	// ignore all tmp files
	if utils.IgnoreFile(event.Name) {
		return true
	}

	return false
}

func (watcher *WSWatcher) ProcessEvent(event *fsnotify.Event) {
	if watcher.ignoreEvent(event) {
		return
	}

	switch {
	// Process create event
	case event.Op&fsnotify.Create == fsnotify.Create:

		// Extract file info
		info, err := os.Stat(event.Name)
		if err == nil {

			if info.IsDir() { // If new item is a directory, track it
				watcher.Add(event.Name)

			} else { // Send event in queue
				watcher.bufferManager.Write(event)
			}
		}

	// Process write event
	case event.Op&fsnotify.Write == fsnotify.Write:
		watcher.bufferManager.Write(event)

	// TODO: Whate about file renaming
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		break
	}
}

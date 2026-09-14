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

func (watcher *WSWatcher) Run() {
	fmt.Println("starting watcher")
	defer func(){
		if err := watcher.Close(); err != nil {
			log.Printf("unable to close filesystem watcher: %v", err)
		}
	}()

	for {
		select {
		case <-watcher.AppContext.Context.Done():
			fmt.Println("watcher stopped properly")
			return
		case event, open := <-watcher.Events:
			if !open {
				log.Println("watcher events channel closed")
				return
			}

			// Process event
			err := watcher.ProcessEvent(&event)
			if err != nil {
				log.Printf("error while processing event: %v\n", err)
			}

		case err, open := <-watcher.Errors:
			if !open {
				log.Println("watcher errors channel closed")
				return
			}

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

func (watcher *WSWatcher) ProcessEvent(event *fsnotify.Event) error {
	if watcher.ignoreEvent(event) {
		return nil
	}

	switch {
	// Process create event
	case event.Op&fsnotify.Create == fsnotify.Create:

		// Extract file info
		info, err := os.Stat(event.Name)
		if err != nil {
			return fmt.Errorf("unable to get file info for %s: %w", event.Name, err)
		} else {
			if info.IsDir() { // If new item is a directory, track it
				if err := watcher.Add(event.Name); err != nil {
					return fmt.Errorf("unable to watch directory %s: %w", event.Name, err)
				}

			} else { // Send event in queue
				watcher.bufferManager.Write(event)
			}
		}

	// Process write event
	case event.Op&fsnotify.Write == fsnotify.Write:
		watcher.bufferManager.Write(event)

	// TODO: What about file renaming
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		break
	}

	return nil
}

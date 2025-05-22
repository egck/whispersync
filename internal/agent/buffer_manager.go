package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/egck/whispersync/internal/types"
	"github.com/fsnotify/fsnotify"
)

type BufferManager struct {
	lock           sync.Mutex
	FlushPeriod    int
	AppContext     *types.AppContext
	requestManager *RequestManager
	SingleBuffers  map[string]*SingleBuffer
}

type SingleBuffer struct {
	event *fsnotify.Event
}

func NewBufferManager(flushPeriod int, AppContext *types.AppContext, requestManager *RequestManager) *BufferManager {
	return &BufferManager{
		FlushPeriod:    flushPeriod,
		AppContext:     AppContext,
		requestManager: requestManager,
		SingleBuffers:  make(map[string]*SingleBuffer),
	}
}

func (bufferManager *BufferManager) Run() {
	fmt.Println("starting buffer manager")
	ticker := time.NewTicker(time.Duration(bufferManager.FlushPeriod) * time.Second)
	for {
		select {
		case <-bufferManager.AppContext.Context.Done():
			fmt.Println("buffer manager stopped properly")
			return
		case <-ticker.C:
			// Flush single buffers
			bufferManager.Flush()
		}
	}
}

func (bufferManager *BufferManager) Flush() {
	// Create slice of events to process
	events := []*fsnotify.Event{}

	// Fill events to process and release lock to manage syncing out of buffers iterations
	bufferManager.lock.Lock()
	for eventName := range bufferManager.SingleBuffers {
		event := bufferManager.SingleBuffers[eventName].extract()

		// Delete old single buffer
		if event == nil {
			delete(bufferManager.SingleBuffers, eventName)
		} else {
			events = append(events, event)
		}
	}
	defer bufferManager.lock.Unlock()

	// Process events
	for _, event := range events {
		bufferManager.requestManager.Sync(event.Name)
	}
}

func (bufferManager *BufferManager) Write(event *fsnotify.Event) {
	bufferManager.lock.Lock()
	defer bufferManager.lock.Unlock()

	if _, exists := bufferManager.SingleBuffers[event.Name]; !exists {
		bufferManager.SingleBuffers[event.Name] = &SingleBuffer{}
	}
	bufferManager.SingleBuffers[event.Name].event = event
}

func (singleBuffer *SingleBuffer) extract() *fsnotify.Event {
	event := singleBuffer.event
	singleBuffer.event = nil

	return event
}

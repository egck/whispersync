package types

import (
	"context"
	"sync"
)

type SyncRequest struct {
	FilePath       string `json:"file_path" binding:"required"`
}

type FileSyncRequest struct {
	SyncRequest
	FileContent []byte `json:"file_content" binding:"required"`
}

type AppContext struct {
	Context   context.Context
	Stop      context.CancelFunc
	WaitGroup sync.WaitGroup
}

func NewAppContext() *AppContext {
	context, cancel := context.WithCancel(context.Background())
	return &AppContext{
		Context: context,
		Stop:    cancel,
	}
}

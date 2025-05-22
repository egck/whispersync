package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/bxgstudio/goconfigloader/loader"
	"github.com/egck/whispersync/internal/config"
	"github.com/egck/whispersync/internal/types"
)

type Core struct {
	Watcher       *WSWatcher
	BufferManager *BufferManager
	AppContext    *types.AppContext
	AppConfig     *config.AgentConfig
}

func NewCore() *Core {
	// Load AppConfig
	AppConfig := &config.AgentConfig{}
	appConfigPath := os.Getenv("WS_AGENT_APP_CONFIG_PATH")
	if appConfigPath == "" {
		appConfigPath = "./app_config.yaml"
	}
	if err := loader.LoadConfig(appConfigPath, AppConfig); err != nil {
		log.Fatalf("unable to load appconfig: %s", err)
	}

	// validate config
	if err := AppConfig.Validate(); err != nil {
		log.Fatalf("invalid appconfig: %s", err)
	}

	// Create AppContext
	AppContext := types.NewAppContext()

	// Create request manager
	requestManager := NewRequestManager(AppConfig.ServerURI, AppConfig.Secret)

	// Create buffer manager
	bufferManager := NewBufferManager(AppConfig.BufferFlushPeriod, AppContext, requestManager)

	// Create watcher
	watcher, err := NewWatcher(
		AppContext,
		AppConfig.WatchDirectoryPath,
		bufferManager,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create GUI
	core := &Core{
		Watcher:       watcher,
		BufferManager: bufferManager,
		AppContext:    AppContext,
		AppConfig:     AppConfig,
	}

	return core
}

func (core *Core) Run() {
	// Run du buffer manager
	core.AppContext.WaitGroup.Add(1)
	go func() {
		defer core.AppContext.WaitGroup.Done()
		core.BufferManager.Run()
	}()

	// Run du watcher
	core.AppContext.WaitGroup.Add(1)
	go func() {
		defer core.AppContext.WaitGroup.Done()
		core.Watcher.Run()
	}()
}

func (core *Core) Stop() {
	fmt.Println("stopping app...")
	core.AppContext.Stop()
	core.AppContext.WaitGroup.Wait()
	fmt.Println("application stopped properly")
}

// TODO: rework
type GetListResponse struct {
	FileNames []string `json:"file_names"`
}

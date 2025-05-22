package server

import (
	"fmt"
	"net/http"

	"github.com/egck/whispersync/internal/types"
	"github.com/gin-gonic/gin"
)

type API struct {
	*http.Server
	Core *Core
}

func NewAPI(core *Core) *API {
	// Init gin router
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// Register routes
	registerRoutes(router, core)

	// Init http server
	httpServer := &http.Server{
		Addr:    core.AppConfig.APIEndpoint,
		Handler: router,
	}

	// Create API
	return &API{Server: httpServer, Core: core}
}

func (api *API) Run() {
	if err := api.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("error from http server: ", err)
	}
}

func (api *API) Stop() {
	if err := api.Shutdown(api.Core.AppContext.Context); err != nil {
		fmt.Println("error while shutting down http server: ", err)
	} else {
		fmt.Println("http server stopped gracefully")
	}
}

func registerRoutes(router *gin.Engine, core *Core) {
	// get status
	router.GET("/api/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// sync file
	router.POST("/api/sync/file", func(context *gin.Context) {
		// Extract json from post request
		var postParams types.FileSyncRequest
		if err := context.ShouldBindJSON(&postParams); err != nil {
			context.JSON(400, gin.H{"message": "bad params"})
			return
		}

		// Try to sync file
		err := core.SyncFile(&postParams)
		if err != nil {
			context.JSON(409, gin.H{"message": err})
			return
		}

		context.JSON(200, gin.H{"message": "sync file ok"})
	})

	// sync rename
	router.POST("/api/sync/rename", func(context *gin.Context) {
		// extract json from sync request
		var postParams types.SyncRequest
		if err := context.ShouldBindJSON(&postParams); err != nil {
			context.JSON(400, gin.H{"message": err})
			return
		}

		// Try to sync rename of the file / directory
		err := core.SyncRename(&postParams)
		if err != nil {
			context.JSON(409, gin.H{"message": err})
			return
		}

		context.JSON(200, gin.H{"message": "sync rename ok"})
	})
}

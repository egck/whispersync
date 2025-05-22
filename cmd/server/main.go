package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/egck/whispersync/internal/server"
)

func main() {
	// Init core
	core := server.NewCore()

	// Run app
	core.Run()

	// Wait for os signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigs

	// Stop app properly
	core.Stop()
	core.AppContext.WaitGroup.Wait()
}

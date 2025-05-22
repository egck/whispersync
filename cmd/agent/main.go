package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/egck/whispersync/internal/agent"
)

func main() {
	// Init & run core
	core := agent.NewCore()
	core.Run()

	// Wait for os signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	// Stop app properly
	core.Stop()
	core.AppContext.WaitGroup.Wait()
}

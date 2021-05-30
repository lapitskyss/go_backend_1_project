package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lapitskyss/go_backend_1_project/internal/di"
)

func main() {
	service, cleanup, err := di.InitializeAPIService()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			service.Log.Error(r.(string))
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	cleanup()
}

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/di"
)

func main() {
	service, cleanup, err := di.InitializeFrontendService()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			service.Log.Error(r.(string))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		service.Log.Infow("Received a signal.", "signal", x.String())
	case err := <-service.Srv.Notify():
		service.Log.Errorw("Received an error from the frontend http server.", "err", err)
	}

	cleanup()
}

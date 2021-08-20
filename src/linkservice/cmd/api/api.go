package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/di"
)

func main() {
	service, cleanup, err := di.InitializeLinkService()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			service.Log.Error(r)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		service.Log.Infow("Received a signal.", "signal", x.String())
	case err := <-service.HTTPServer.Notify():
		service.Log.Errorw("Received an error from the linkservice http server.", "err", err)
	}

	cleanup()
}

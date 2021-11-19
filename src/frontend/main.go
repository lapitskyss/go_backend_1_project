package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

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
		service.Log.Info("Frontend received a signal.", zap.String("signal", x.String()))
	case e := <-service.Srv.Notify():
		service.Log.Error("Received an error from the frontend http server.", zap.Error(e))
	}

	cleanup()
}

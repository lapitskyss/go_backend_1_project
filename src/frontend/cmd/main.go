package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/di"
)

func main() {
	f, cleanup, err := di.InitializeFrontend()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			f.Log.Error("panic", zap.Any("details", r))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		f.Log.Info("Frontend received a signal.", zap.String("signal", x.String()))
	case e := <-f.Srv.Notify():
		f.Log.Error("Received an error from the frontend http server.", zap.Error(e))
	}

	cleanup()
}

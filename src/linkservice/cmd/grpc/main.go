package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/di"
)

func main() {
	service, cleanup, err := di.InitializeGRPC()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			service.Log.Error("panic", zap.Any("details", r))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		service.Log.Info("Received a signal. Stopped link grpc service.", zap.String("signal", x.String()))
	case e := <-service.GRPCServer.Notify():
		service.Log.Error("Received an error from the linkservice grpc server.", zap.Error(e))
	}

	cleanup()
}

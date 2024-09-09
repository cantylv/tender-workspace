package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"tender-workspace/internal/delivery/route"
	"tender-workspace/services/postgres"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger) {
	psqlConn := postgres.Init(logger)

	r := mux.NewRouter()
	handler := route.InitHTTPHandlers(r, psqlConn, logger)

	srv := &http.Server{
		Handler:      handler,
		Addr:         viper.GetString("SERVER_ADDRESS"),
		WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
		IdleTimeout:  viper.GetDuration("SERVER_IDLE_TIMEOUT"),
	}

	go func() {
		logger.Info(fmt.Sprintf("server has started at the address %s", viper.GetString("SERVER_ADDRESS")))
		if err := srv.ListenAndServe(); err != nil {
			logger.Warn(fmt.Sprintf("error after end of receiving requests: %v", err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("SERVER_SHUTDOWN_DURATION"))
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("server urgently has shut down with an error: %v", err))
		os.Exit(1)
	}
	logger.Info("server has shut down")
	os.Exit(0)
}

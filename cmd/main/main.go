package main

import (
	"tender-workspace/config"
	"tender-workspace/internal/app"

	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	config.Read("./config/config.yaml", logger)
	app.Run(logger)
}

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Read(path string, logger *zap.Logger) {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logger.Fatal(fmt.Sprintf("fatal error config file: %v", err))
		}
		logger.Error(fmt.Sprintf("error while reading conf file: %v", err))
		logger.Info("configuration file is not found, programm will be executed within default configuration")
	}
	logger.Info("successful read of configuration")
}

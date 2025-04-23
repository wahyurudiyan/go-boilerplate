package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func getEnvFile() string {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		slog.Warn(".env not found, falling back to .env.example")
		envFile = ".env.example"
	}

	return envFile
}

func Load[T any](out *T) {
	file := getEnvFile()

	viper.SetConfigFile(file)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&out); err != nil {
		panic(err)
	}
	slog.Info("Configuration loaded from", "filename", file)
}

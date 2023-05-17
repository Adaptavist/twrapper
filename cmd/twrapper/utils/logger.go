package utils

import (
	"os"

	"go.uber.org/zap"
)

func GetLogger() *zap.Logger {
	if os.Getenv("TW_ENV") == "development" {
		return zap.Must(zap.NewDevelopment())
	}

	return zap.Must(zap.NewProduction())
}

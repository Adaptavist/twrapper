// Providing helpers for our integration
package helpers

import (
	"os"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"go.uber.org/zap"
)

func Config() *config.Config {
	return config.Must(
		config.NewTest(
			zap.Must(zap.NewDevelopment())))
}

func CleanupDIR(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
}

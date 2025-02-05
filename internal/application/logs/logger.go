package logs

import (
	"log/slog"
	"os"
	"strings"

	"github.com/openinnovationai/k8s-amd-exporter/internal/application/settings"
)

const developmentLog string = "development"

func MakeLogger(configuration *settings.Configuration) *slog.Logger {
	logLevel := slog.LevelInfo

	if strings.ToLower(configuration.LogLevel) == developmentLog {
		logLevel = slog.LevelDebug
	}

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	slog.SetLogLoggerLevel(logLevel)

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	logger := slog.New(loggerHandler)

	logger.Info("log settings", slog.String("level", handlerOptions.Level.Level().String()))

	return logger
}

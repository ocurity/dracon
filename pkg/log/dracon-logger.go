package draconLogger

import (
	"log/slog"
	"os"
)

func SetDefault(logLevel slog.Level, scanID string, jsonLogging bool) {
	var logger *slog.Logger
	if jsonLogging {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}
	if scanID != "" {
		logger = logger.With("scanID", scanID)
	}
	slog.SetDefault(logger)
}

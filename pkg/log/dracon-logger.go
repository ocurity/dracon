package draconLogger

import (
	"log/slog"
	"os"
)

func SetDefault(logLevel slog.Level) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
}

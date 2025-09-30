package logger

import (
	"io"
	"log/slog"
)

func SetupLogger(logWriter io.Writer, logLevel string) {
	options := new(slog.HandlerOptions)
	switch logLevel {
	case "debug":
		options.Level = slog.LevelDebug
	case "info":
		options.Level = slog.LevelInfo
	case "warn":
		options.Level = slog.LevelWarn
	case "error":
		options.Level = slog.LevelError
	default:
		options.Level = slog.LevelInfo
		slog.Warn("Incorrect logger level. The level is set to info.")
	}
	logHandler := slog.NewTextHandler(logWriter, options)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

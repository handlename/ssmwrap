package ssmwrap

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func InitLogger() {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     selectLogLevel(),
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &opts))
	slog.SetDefault(logger)
}

func selectLogLevel() slog.Leveler {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "" {
		level = "info"
	}

	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		log.Println("invalid log level, using info level.")
		return slog.LevelInfo
	}
}

package pkg

import (
	"io"
	"log/slog"
)

func CreateLogger(logLevel slog.Level, logDestination io.Writer) *slog.Logger {
	logOptions := slog.HandlerOptions{Level: logLevel}
	logHandler := slog.NewJSONHandler(logDestination, &logOptions)

	return slog.New(logHandler)
}

func ParseLevel(strLevel string) slog.Level {
	var level slog.Level

	var err = level.UnmarshalText([]byte(strLevel))
	if err != nil {
		level = slog.LevelInfo
	}

	return level
}

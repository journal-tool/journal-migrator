package test

import (
	"bytes"
	"journal-migrator/pkg"
	"journal-migrator/pkg/storage"
	"log/slog"
	"os"
)

func SetupDatabaseClient(dbName string) *storage.DatabaseClient {
	config := storage.NewDatabaseConfig(
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASS"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_TYPE"),
	)

	if config.Name == "" {
		config.Name = dbName
	}

	return storage.NewDatabaseClient(config)
}

func SetupLogger() *slog.Logger {
	memoryBuffer := bytes.Buffer{}
	memoryLogger := pkg.CreateLogger(slog.LevelInfo, &memoryBuffer)

	return memoryLogger
}

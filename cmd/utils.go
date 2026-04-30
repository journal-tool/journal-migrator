package cmd

import (
	"github.com/urfave/cli/v3"
	"journal-migrator/pkg"
	"journal-migrator/pkg/storage"
	"log/slog"
	"os"
)

func SetupDatabaseClient(command *cli.Command) *storage.DatabaseClient {
	config := storage.NewDatabaseConfig(
		command.String("databaseUsername"),
		command.String("databasePassword"),
		command.String("databaseHost"),
		command.String("databasePort"),
		command.String("databaseName"),
		command.String("databaseType"),
	)

	return storage.NewDatabaseClient(config)
}

func SetupLogger(command *cli.Command) *slog.Logger {
	logLevel := command.String("logLevel")

	level := pkg.ParseLevel(logLevel)
	logger := pkg.CreateLogger(level, os.Stdout)

	return logger
}

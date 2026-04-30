package main

import (
	"context"
	"github.com/urfave/cli/v3"
	"journal-migrator/cmd"
	"journal-migrator/pkg/routines"
	"log"
	"log/slog"
	"os"
)

func run(_ context.Context, command *cli.Command) error {
	logger := cmd.SetupLogger(command)
	client := cmd.SetupDatabaseClient(command)

	err := client.OpenConnection()
	if err != nil {
		return err
	}

	defer client.CloseConnection()

	table := command.String("table")
	strategy := command.String("strategy")

	err = routines.NewCleanupRoutine(client, logger, strategy).Run(table)
	if err != nil {
		logger.Error("Routine failed",
			slog.String("reason", err.Error()),
			slog.String("table", table),
			slog.String("strategy", strategy),
		)
	}

	return err
}

func main() {
	app := &cli.Command{
		Name: "Journal Migrator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "databaseUsername",
				Value:    "unknown",
				Usage:    "Database connection username",
				Required: true,
				Sources:  cli.EnvVars("DATABASE_USER"),
			},
			&cli.StringFlag{
				Name:     "databasePassword",
				Value:    "unknown",
				Usage:    "Database connection password",
				Required: true,
				Sources:  cli.EnvVars("DATABASE_PASS"),
			},
			&cli.StringFlag{
				Name:     "databaseHost",
				Value:    "unknown",
				Usage:    "Database connection host",
				Required: true,
				Sources:  cli.EnvVars("DATABASE_HOST"),
			},
			&cli.StringFlag{
				Name:     "databasePort",
				Value:    "unknown",
				Usage:    "Database connection port",
				Required: true,
				Sources:  cli.EnvVars("DATABASE_PORT"),
			},
			&cli.StringFlag{
				Name:     "databaseName",
				Value:    "unknown",
				Usage:    "Database connection name",
				Required: true,
				Sources:  cli.EnvVars("DATABASE_NAME"),
			},
			&cli.StringFlag{
				Name:     "databaseType",
				Value:    "unknown",
				Usage:    "Database type",
				Required: true,
				Sources:  cli.EnvVars("DATABASE_TYPE"),
			},
			&cli.StringFlag{
				Name:     "logLevel",
				Value:    "INFO",
				Usage:    "Logger severity level",
				Required: false,
				Sources:  cli.EnvVars("LOG_LEVEL"),
			},
			&cli.StringFlag{
				Name:     "table",
				Value:    "unknown",
				Usage:    "Cleanup table to target",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "strategy",
				Value:    routines.SyncStrategy,
				Usage:    "Cleanup strategy",
				Required: false,
			},
		},
		Action: run,
	}

	ctx := context.Background()

	err := app.Run(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

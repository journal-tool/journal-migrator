package routines

import (
	"journal-migrator/pkg/routines/ctx"
	"journal-migrator/pkg/routines/routines/cleanup"
	"journal-migrator/pkg/routines/routines/migrate"
	"journal-migrator/pkg/storage"
	"log/slog"
)

const AsyncStrategy = "ASYNC"
const SyncStrategy = "SYNC"

type BaseCleanupRoutine = cleanup.BaseCleanupRoutine
type BaseMigrateRoutine = migrate.BaseMigrateRoutine

func NewCleanupRoutine(client *storage.DatabaseClient, logger *slog.Logger, strategy string) cleanup.BaseCleanupRoutine {
	config := client.GetConfig()
	context := ctx.NewContext(config, logger)

	switch strategy {
	case SyncStrategy:
		return cleanup.NewSyncCleanupRoutine(context, client, logger)
	}

	return nil
}

func NewMigrateRoutine(client *storage.DatabaseClient, logger *slog.Logger, strategy string) migrate.BaseMigrateRoutine {
	config := client.GetConfig()
	context := ctx.NewContext(config, logger)

	switch strategy {
	case AsyncStrategy:
		return migrate.NewAsyncMigrateRoutine(context, client, logger)
	case SyncStrategy:
		return migrate.NewSyncMigrateRoutine(context, client, logger)
	}

	return nil
}

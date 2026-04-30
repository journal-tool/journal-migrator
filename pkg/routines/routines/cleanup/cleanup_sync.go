package cleanup

import (
	"fmt"
	"journal-migrator/pkg/routines/ctx"
	"journal-migrator/pkg/storage"
	"log/slog"
	"sync/atomic"
)

type SyncCleanupRoutine struct {
	context  *ctx.Context
	client   *storage.DatabaseClient
	logger   *slog.Logger
	progress atomic.Int64
}

func NewSyncCleanupRoutine(context *ctx.Context, client *storage.DatabaseClient, logger *slog.Logger) *SyncCleanupRoutine {
	return &SyncCleanupRoutine{
		context: context,
		client:  client,
		logger:  logger,
	}
}

func (r *SyncCleanupRoutine) Progress() int64 {
	return r.progress.Load()
}

func (r *SyncCleanupRoutine) Run(table string) error {
	sourceTable := table
	targetTable := fmt.Sprintf("_migrator_%s", table)
	legacyTable := fmt.Sprintf("_archived_%s", table)

	database := r.client.GetDatabase()
	executor := r.context.CreateExecutor(database)

	err := executor.DeleteTableTriggers(sourceTable)
	if err != nil {
		return err
	}

	err = executor.DeleteTable(targetTable)
	if err != nil {
		return err
	}

	err = executor.DeleteTable(legacyTable)
	if err != nil {
		return err
	}

	r.progress.Store(100)
	return nil
}

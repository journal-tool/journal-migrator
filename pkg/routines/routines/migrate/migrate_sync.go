package migrate

import (
	"errors"
	"journal-migrator/pkg/handlers/throttlers"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/routines/ctx"
	"journal-migrator/pkg/storage"
	"log/slog"
	"sync/atomic"
)

type SyncMigrateRoutine struct {
	context  *ctx.Context
	client   *storage.DatabaseClient
	logger   *slog.Logger
	progress atomic.Int64
}

func NewSyncMigrateRoutine(context *ctx.Context, client *storage.DatabaseClient, logger *slog.Logger) *SyncMigrateRoutine {
	return &SyncMigrateRoutine{
		context: context,
		client:  client,
		logger:  logger,
	}
}

func (r *SyncMigrateRoutine) migrateTable(table string, operations []models.Operation) error {
	database := r.client.GetDatabase()
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	migrator := r.context.CreateMigrator(tx)

	err = migrator.Migrate(table, operations)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SyncMigrateRoutine) validateOperations(operations []models.Operation) error {
	if len(operations) == 0 {
		return errors.New("cannot perform zero table operations")
	}

	for _, operation := range operations {
		if operation.IsTableOperation() && len(operations) > 1 {
			return errors.New("cannot perform table operations in bulk")
		}
	}

	return nil
}

func (r *SyncMigrateRoutine) Progress() int64 {
	return r.progress.Load()
}

func (r *SyncMigrateRoutine) Run(table string, operations []models.Operation, _ throttlers.BaseThrottler) error {
	err := r.validateOperations(operations)
	if err != nil {
		return err
	}

	err = r.migrateTable(table, operations)
	if err != nil {
		return err
	}

	r.progress.Store(100)
	return nil
}

package migrate

import (
	"errors"
	"fmt"
	"journal-migrator/pkg/handlers/throttlers"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/routines/ctx"
	"journal-migrator/pkg/storage"
	"log/slog"
	"sync/atomic"
)

type AsyncMigrateRoutine struct {
	context   *ctx.Context
	client    *storage.DatabaseClient
	logger    *slog.Logger
	progress  atomic.Int64
	batchSize int
}

func NewAsyncMigrateRoutine(context *ctx.Context, client *storage.DatabaseClient, logger *slog.Logger) *AsyncMigrateRoutine {
	return &AsyncMigrateRoutine{
		context:   context,
		client:    client,
		logger:    logger,
		batchSize: 1000,
	}
}

func (r *AsyncMigrateRoutine) backfillTable(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string, throttler throttlers.BaseThrottler) error {
	database := r.client.GetDatabase()
	executor := r.context.CreateExecutor(database)
	fetcher := r.context.CreateFetcher(database)

	lowerBoundID := 0
	upperBoundID := r.batchSize

	maxID, err := fetcher.FetchColumnMax(sourceTable.Name, sourceTable.Key)
	if err != nil {
		return err
	}

	for lowerBoundID < maxID {
		err = executor.InsertRows(sourceTable, targetTable, columns, lowerBoundID, upperBoundID)
		if err != nil {
			return err
		}

		progressNow := min(upperBoundID, maxID)
		progressRatio := float64(progressNow) / float64(maxID)
		r.progress.Store(int64(progressRatio * 100))

		maxID, err = fetcher.FetchColumnMax(sourceTable.Name, sourceTable.Key)
		if err != nil {
			return err
		}

		throttler.Throttle()
		lowerBoundID += r.batchSize
		upperBoundID += r.batchSize
	}

	return nil
}

func (r *AsyncMigrateRoutine) entangleTable(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string, wrappedFn func() error) error {
	database := r.client.GetDatabase()
	executor := r.context.CreateExecutor(database)

	err := executor.CreateTableTriggers(sourceTable, targetTable, columns)
	if err != nil {
		return err
	}

	wrappedError := wrappedFn()

	err = executor.DeleteTableTriggers(sourceTable.Name)
	if err != nil {
		return err
	}

	return wrappedError
}

func (r *AsyncMigrateRoutine) migrateTable(table string, operations []models.Operation) error {
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

func (r *AsyncMigrateRoutine) renameTables(sourceTable string, targetTable string, legacyTable string) error {
	database := r.client.GetDatabase()
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	migrator := r.context.CreateMigrator(tx)

	err = migrator.RenameTable(sourceTable, legacyTable)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = migrator.RenameTable(targetTable, sourceTable)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *AsyncMigrateRoutine) validateOperations(operations []models.Operation) error {
	if len(operations) == 0 {
		return errors.New("cannot perform zero table operations")
	}

	for _, operation := range operations {
		if operation.IsTableOperation() {
			return errors.New("cannot perform table operations asynchronously")
		}
	}

	return nil
}

func (r *AsyncMigrateRoutine) validateTableKey(table string) error {
	database := r.client.GetDatabase()
	fetcher := r.context.CreateFetcher(database)

	info, err := fetcher.FetchTableInfo(table)
	if err != nil {
		return err
	}
	if info.Valid != true {
		return errors.New("table does not have a valid auto-incremented key")
	}

	return nil
}

func (r *AsyncMigrateRoutine) Progress() int64 {
	return r.progress.Load()
}

func (r *AsyncMigrateRoutine) Run(table string, operations []models.Operation, throttler throttlers.BaseThrottler) error {
	err := r.validateOperations(operations)
	if err != nil {
		return err
	}

	err = r.validateTableKey(table)
	if err != nil {
		return err
	}

	sourceTable := table
	targetTable := fmt.Sprintf("_migrator_%s", table)
	legacyTable := fmt.Sprintf("_archived_%s", table)

	database := r.client.GetDatabase()
	executor := r.context.CreateExecutor(database)
	fetcher := r.context.CreateFetcher(database)

	err = executor.CreateTableDuplicate(sourceTable, targetTable)
	if err != nil {
		return err
	}

	err = r.migrateTable(targetTable, operations)
	if err != nil {
		return err
	}

	sourceTableInfo, _ := fetcher.FetchTableInfo(sourceTable)
	targetTableInfo, _ := fetcher.FetchTableInfo(targetTable)
	sharedTableCols, _ := fetcher.FetchTableColumnIntersection(sourceTable, targetTable)

	err = r.entangleTable(*sourceTableInfo, *targetTableInfo, sharedTableCols, func() error {
		return r.backfillTable(*sourceTableInfo, *targetTableInfo, sharedTableCols, throttler)
	})
	if err != nil {
		return err
	}

	err = r.renameTables(sourceTable, targetTable, legacyTable)
	if err != nil {
		return err
	}

	return nil
}

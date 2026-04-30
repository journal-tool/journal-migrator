package handlers

import (
	"fmt"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"log/slog"
	"strings"
)

type ExecutorHandler struct {
	dialector dialectors.BaseExecutorDialector
	database  storage.Database
	config    storage.DatabaseConfig
	logger    *slog.Logger
}

func NewExecutorHandler(dialector dialectors.BaseExecutorDialector, database storage.Database, config storage.DatabaseConfig, logger *slog.Logger) *ExecutorHandler {
	return &ExecutorHandler{
		dialector: dialector,
		database:  database,
		config:    config,
		logger:    logger,
	}
}

func (h *ExecutorHandler) InsertRows(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string, lowerBound int, upperBound int) error {
	var query = h.dialector.TableInsertBatchQuery(sourceTable, targetTable, columns, lowerBound, upperBound)

	result, err := h.database.Exec(query)
	if err != nil {
		h.logger.Error("Cannot insert rows",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", sourceTable.Name),
		)
		return err
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		h.logger.Error("Inserted rows unknown",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", sourceTable.Name),
		)
	} else {
		h.logger.Debug(
			fmt.Sprintf("Inserted rows: %d", numRows),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", sourceTable.Name),
		)
	}

	return err
}

func (h *ExecutorHandler) CreateTableDuplicate(sourceTable string, targetTable string) error {
	var query = h.dialector.TableCreateQuery(sourceTable, targetTable)

	_, err := h.database.Exec(query)
	if err != nil {
		h.logger.Error("Cannot create table duplicate",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", sourceTable),
		)
		return err
	}

	return nil
}

func (h *ExecutorHandler) DeleteTable(legacyTable string) error {
	var query = h.dialector.TableRemoveQuery(legacyTable)

	_, err := h.database.Exec(query)
	if err != nil {
		h.logger.Error("Cannot delete table",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", legacyTable),
		)
		return err
	}

	return nil
}

func (h *ExecutorHandler) CreateTableTriggers(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string) error {
	var queries = []string{
		h.dialector.TriggerCreateDeleteQuery(sourceTable, targetTable),
		h.dialector.TriggerCreateUpdateQuery(sourceTable, targetTable, columns),
		h.dialector.TriggerCreateInsertQuery(sourceTable, targetTable, columns),
	}

	var query = strings.Join(queries, ";")

	_, err := h.database.Exec(query)
	if err != nil {
		h.logger.Error("Cannot create table triggers",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", sourceTable.Name),
		)
		return err
	}

	return nil
}

func (h *ExecutorHandler) DeleteTableTriggers(sourceTable string) error {
	var queries = []string{
		h.dialector.TriggerDropDeleteQuery(sourceTable),
		h.dialector.TriggerDropUpdateQuery(sourceTable),
		h.dialector.TriggerDropInsertQuery(sourceTable),
	}

	var query = strings.Join(queries, ";")

	_, err := h.database.Exec(query)
	if err != nil {
		h.logger.Error("Cannot drop table triggers",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", sourceTable),
		)
		return err
	}

	return nil
}

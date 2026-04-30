package handlers

import (
	"errors"
	"journal-migrator/lib"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"log/slog"
)

var ParseColumnInfoRows = lib.ParseRows[models.ColumnInfo]
var ParseIndexInfoRows = lib.ParseRows[models.IndexInfo]
var ParseTableInfoRow = lib.ParseRow[models.TableInfo]
var ParseTableSizeRow = lib.ParseRow[models.TableSize]

type FetcherHandler struct {
	dialector dialectors.BaseFetcherDialector
	database  storage.Database
	config    storage.DatabaseConfig
	logger    *slog.Logger
}

func NewFetcherHandler(dialector dialectors.BaseFetcherDialector, database storage.Database, config storage.DatabaseConfig, logger *slog.Logger) *FetcherHandler {
	return &FetcherHandler{
		dialector: dialector,
		database:  database,
		config:    config,
		logger:    logger,
	}
}

func (h *FetcherHandler) FetchColumnMax(table string, column string) (int, error) {
	var query = h.dialector.ColumnSelectMaxQuery(table, column)
	var maxVal int

	row := h.database.QueryRow(query)
	err := row.Scan(&maxVal)
	if err != nil {
		h.logger.Error("Cannot query table maximum key",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", table),
		)
		return maxVal, err
	}

	return maxVal, nil
}

func (h *FetcherHandler) FetchTableInfo(table string) (*models.TableInfo, error) {
	var query = h.dialector.TableSelectInfoQuery(table)

	row := h.database.QueryRow(query)
	err := row.Err()
	if err != nil {
		h.logger.Error("Cannot query table information",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", table),
		)
		return nil, err
	}

	return ParseTableInfoRow(row)
}

func (h *FetcherHandler) FetchTableSize(table string) (*models.TableSize, error) {
	var query = h.dialector.TableSelectSizeQuery(table)

	row := h.database.QueryRow(query)
	err := row.Err()
	if err != nil {
		h.logger.Error("Cannot query table size",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", table),
		)
		return nil, err
	}

	return ParseTableSizeRow(row)
}

func (h *FetcherHandler) FetchTableColumns(table string) ([]models.ColumnInfo, error) {
	var query = h.dialector.TableSelectColumnsQuery(table)

	rows, err := h.database.Query(query)
	if err != nil {
		h.logger.Error("Cannot query table column information",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", table),
		)
		return nil, err
	}

	return ParseColumnInfoRows(rows)
}

func (h *FetcherHandler) FetchTableIndexes(table string) ([]models.IndexInfo, error) {
	var query = h.dialector.TableSelectIndexesQuery(table)

	rows, err := h.database.Query(query)
	if err != nil {
		h.logger.Error("Cannot query table index information",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", table),
		)
		return nil, err
	}

	return ParseIndexInfoRows(rows)
}

func (h *FetcherHandler) FetchTableColumnIntersection(sourceTable string, targetTable string) ([]string, error) {
	sourceTableColumns, err := h.FetchTableColumns(sourceTable)
	if err != nil {
		return nil, err
	}

	targetTableColumns, err := h.FetchTableColumns(targetTable)
	if err != nil {
		return nil, err
	}

	commonColumns := models.IntersectColumns(sourceTableColumns, targetTableColumns)
	if len(commonColumns) == 0 {
		return nil, errors.New("no column intersection")
	}

	return commonColumns, nil
}

func (h *FetcherHandler) FetchTableIndexIntersection(sourceTable string, targetTable string) ([]string, error) {
	sourceTableIndexes, err := h.FetchTableIndexes(sourceTable)
	if err != nil {
		return nil, err
	}

	targetTableIndexes, err := h.FetchTableIndexes(targetTable)
	if err != nil {
		return nil, err
	}

	commonIndexes := models.IntersectIndexes(sourceTableIndexes, targetTableIndexes)
	if len(commonIndexes) == 0 {
		return nil, errors.New("no index intersection")
	}

	return commonIndexes, nil
}

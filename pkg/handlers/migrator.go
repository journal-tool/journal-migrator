package handlers

import (
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/models/operation"
	"journal-migrator/pkg/models/operation/specs"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"log/slog"
	"strings"
)

type MigratorHandler struct {
	dialector dialectors.BaseMigratorDialector
	database  storage.Database
	config    storage.DatabaseConfig
	logger    *slog.Logger
}

func NewMigratorHandler(dialector dialectors.BaseMigratorDialector, database storage.Database, config storage.DatabaseConfig, logger *slog.Logger) *MigratorHandler {
	return &MigratorHandler{
		dialector: dialector,
		database:  database,
		config:    config,
		logger:    logger,
	}
}

func (h *MigratorHandler) buildQuery(table string, op models.Operation) string {
	switch op.Type {
	case operation.CreateColumnType:
		spec := op.Spec.(*specs.CreateColumnSpec)
		return h.dialector.CreateColumnQuery(
			table,
			spec.Name,
			spec.Type,
			spec.Default,
			spec.Nullable,
		)
	case operation.ChangeColumnType:
		spec := op.Spec.(*specs.ChangeColumnSpec)
		return h.dialector.ChangeColumnQuery(
			table,
			spec.Name,
			spec.Type,
			spec.Default,
			spec.Nullable,
		)
	case operation.RemoveColumnType:
		spec := op.Spec.(*specs.RemoveColumnSpec)
		return h.dialector.RemoveColumnQuery(
			table,
			spec.Name,
		)
	case operation.RenameColumnType:
		spec := op.Spec.(*specs.RenameColumnSpec)
		return h.dialector.RenameColumnQuery(
			table,
			spec.OldName,
			spec.NewName,
		)
	case operation.CreateIndexType:
		spec := op.Spec.(*specs.CreateIndexSpec)
		return h.dialector.CreateIndexQuery(
			table,
			spec.Name,
			spec.Columns,
			spec.Unique,
		)
	case operation.ChangeIndexType:
		spec := op.Spec.(*specs.ChangeIndexSpec)
		return h.dialector.ChangeIndexQuery(
			table,
			spec.Name,
			spec.Columns,
			spec.Unique,
		)
	case operation.RemoveIndexType:
		spec := op.Spec.(*specs.RemoveIndexSpec)
		return h.dialector.RemoveIndexQuery(
			table,
			spec.Name,
		)
	case operation.RenameIndexType:
		spec := op.Spec.(*specs.RenameIndexSpec)
		return h.dialector.RenameIndexQuery(
			table,
			spec.OldName,
			spec.NewName,
		)
	case operation.CreateTableType:
		spec := op.Spec.(*specs.CreateTableSpec)
		return h.dialector.CreateTableQuery(
			spec.Name,
			spec.Definition,
		)
	case operation.RemoveTableType:
		spec := op.Spec.(*specs.RemoveTableSpec)
		return h.dialector.RemoveTableQuery(
			spec.Name,
		)
	case operation.RenameTableType:
		spec := op.Spec.(*specs.RenameTableSpec)
		return h.dialector.RenameTableQuery(
			spec.OldName,
			spec.NewName,
		)
	case operation.DDLType:
		spec := op.Spec.(*specs.DDLSpec)
		return h.dialector.DDLOperationQuery(
			spec.Statement,
		)
	}

	return ""
}

func (h *MigratorHandler) Migrate(table string, operations []models.Operation) error {
	var queries []string
	for _, op := range operations {
		queries = append(queries, h.buildQuery(table, op))
	}

	var query = strings.Join(queries, ";")

	_, err := h.database.Exec(query)
	if err != nil {
		h.logger.Error("Cannot apply operations",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", table),
		)
		return err
	}

	return nil
}

func (h *MigratorHandler) RenameTable(oldName string, newName string) error {
	var query = h.dialector.RenameTableQuery(oldName, newName)

	_, err := h.database.Exec(query)
	if err != nil {
		h.logger.Error("Cannot rename tables",
			slog.String("reason", err.Error()),
			slog.String("db_host", h.config.Host),
			slog.String("db_port", h.config.Port),
			slog.String("db_name", h.config.Name),
			slog.String("db_table", oldName),
		)
		return err
	}

	return nil
}

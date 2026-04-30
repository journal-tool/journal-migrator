package executor

import (
	"fmt"
	"journal-migrator/lib"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/storage/formatters"
	"strings"
)

type PostgreSQLExecutorDialector struct {
	formatter *formatters.PostgreSQLFormatter
	schema    string
}

func NewPostgreSQLExecutorDialector(formatter *formatters.PostgreSQLFormatter, schema string) *PostgreSQLExecutorDialector {
	return &PostgreSQLExecutorDialector{
		formatter: formatter,
		schema:    schema,
	}
}

func (d *PostgreSQLExecutorDialector) buildDeleteQuery(targetTable models.TableInfo) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DELETE FROM %s WHERE %s = %s`),
		d.formatter.QuoteEntity(targetTable.Name),
		d.formatter.ScopeEntity(targetTable.Name, targetTable.Key),
		d.scopeOldColumn(targetTable.Key),
	)
}

func (d *PostgreSQLExecutorDialector) buildUpsertQuery(targetTable models.TableInfo, columns []string) string {
	columnsUpdateFunc := func(column string) string { return fmt.Sprintf("%s = %s", column, d.scopeExcColumn(column)) }

	columnsQuoted := lib.Mapper(columns, d.formatter.QuoteEntity)
	columnsScopped := lib.Mapper(columns, d.scopeNewColumn)
	columnsUpdated := lib.Mapper(columns, columnsUpdateFunc)

	return fmt.Sprintf(
		d.formatter.TrimSpaces(`INSERT INTO %s (%s) VALUES (%s) ON CONFLICT(%s) DO UPDATE SET %s`),
		d.formatter.QuoteEntity(targetTable.Name),
		strings.Join(columnsQuoted, ", "),
		strings.Join(columnsScopped, ", "),
		d.formatter.QuoteEntity(targetTable.Key),
		strings.Join(columnsUpdated, ", "),
	)
}

func (d *PostgreSQLExecutorDialector) scopeExcColumn(column string) string {
	return fmt.Sprintf("EXCLUDED.%s", d.formatter.QuoteEntity(column))
}

func (d *PostgreSQLExecutorDialector) scopeNewColumn(column string) string {
	return fmt.Sprintf("NEW.%s", d.formatter.QuoteEntity(column))
}

func (d *PostgreSQLExecutorDialector) scopeOldColumn(column string) string {
	return fmt.Sprintf("OLD.%s", d.formatter.QuoteEntity(column))
}

func (d *PostgreSQLExecutorDialector) TableCreateQuery(sourceTable string, targetTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE TABLE %s (LIKE %s INCLUDING ALL)`),
		d.formatter.ScopeEntity(d.schema, targetTable),
		d.formatter.ScopeEntity(d.schema, sourceTable),
	)
}

func (d *PostgreSQLExecutorDialector) TableRemoveQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP TABLE IF EXISTS %s`),
		d.formatter.ScopeEntity(d.schema, sourceTable),
	)
}

func (d *PostgreSQLExecutorDialector) TableInsertBatchQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string, lowerID int, upperID int) string {
	scopedColumn := d.formatter.ScopeEntity(sourceTable.Name, sourceTable.Key)
	scopedSourceTable := d.formatter.ScopeEntity(d.schema, sourceTable.Name)
	scopedTargetTable := d.formatter.ScopeEntity(d.schema, targetTable.Name)

	columnsQuoted := lib.Mapper(columns, d.formatter.QuoteEntity)
	columnsJoined := strings.Join(columnsQuoted, ", ")

	selectQuery := fmt.Sprintf(
		d.formatter.TrimSpaces(`SELECT %s FROM %s WHERE %s BETWEEN %d AND %d`),
		columnsJoined,
		scopedSourceTable,
		scopedColumn,
		lowerID,
		upperID,
	)

	return fmt.Sprintf(
		d.formatter.TrimSpaces(`INSERT INTO %s (%s) %s ON CONFLICT DO NOTHING`),
		scopedTargetTable,
		columnsJoined,
		selectQuery,
	)
}

func (d *PostgreSQLExecutorDialector) TriggerCreateDeleteQuery(sourceTable models.TableInfo, targetTable models.TableInfo) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE RULE %s AS ON DELETE TO %s DO ALSO (%s)`),
		d.formatter.PrefixName(`_delete`, sourceTable.Name),
		d.formatter.ScopeEntity(d.schema, sourceTable.Name),
		d.buildDeleteQuery(targetTable),
	)
}

func (d *PostgreSQLExecutorDialector) TriggerCreateInsertQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE RULE %s AS ON INSERT TO %s DO ALSO (%s)`),
		d.formatter.PrefixName(`_insert`, sourceTable.Name),
		d.formatter.ScopeEntity(d.schema, sourceTable.Name),
		d.buildUpsertQuery(targetTable, columns),
	)
}

func (d *PostgreSQLExecutorDialector) TriggerCreateUpdateQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE RULE %s AS ON UPDATE TO %s DO ALSO (%s)`),
		d.formatter.PrefixName(`_update`, sourceTable.Name),
		d.formatter.ScopeEntity(d.schema, sourceTable.Name),
		d.buildUpsertQuery(targetTable, columns),
	)
}

func (d *PostgreSQLExecutorDialector) TriggerDropDeleteQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP RULE IF EXISTS %s ON %s`),
		d.formatter.PrefixName(`_delete`, sourceTable),
		d.formatter.ScopeEntity(d.schema, sourceTable),
	)
}

func (d *PostgreSQLExecutorDialector) TriggerDropInsertQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP RULE IF EXISTS %s ON %s`),
		d.formatter.PrefixName(`_insert`, sourceTable),
		d.formatter.ScopeEntity(d.schema, sourceTable),
	)
}

func (d *PostgreSQLExecutorDialector) TriggerDropUpdateQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP RULE IF EXISTS %s ON %s`),
		d.formatter.PrefixName(`_update`, sourceTable),
		d.formatter.ScopeEntity(d.schema, sourceTable),
	)
}

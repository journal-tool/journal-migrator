package executor

import (
	"fmt"
	"journal-migrator/lib"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/storage/formatters"
	"strings"
)

type MySQLExecutorDialector struct {
	formatter *formatters.MySQLFormatter
	schema    string
}

func NewMySQLExecutorDialector(formatter *formatters.MySQLFormatter, schema string) *MySQLExecutorDialector {
	return &MySQLExecutorDialector{
		formatter: formatter,
		schema:    schema,
	}
}

func (d *MySQLExecutorDialector) buildDeleteQuery(targetTable models.TableInfo) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DELETE FROM %s WHERE %s = %s`),
		d.formatter.QuoteEntity(targetTable.Name),
		d.formatter.ScopeEntity(targetTable.Name, targetTable.Key),
		d.scopeOldColumn(targetTable.Key),
	)
}

func (d *MySQLExecutorDialector) buildUpsertQuery(targetTable models.TableInfo, columns []string) string {
	columnsUpdateFunc := func(column string) string { return fmt.Sprintf("%s = %s", column, d.scopeNewColumn(column)) }

	columnsQuoted := lib.Mapper(columns, d.formatter.QuoteEntity)
	columnsScopped := lib.Mapper(columns, d.scopeNewColumn)
	columnsUpdated := lib.Mapper(columns, columnsUpdateFunc)

	return fmt.Sprintf(
		d.formatter.TrimSpaces(`INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s`),
		d.formatter.QuoteEntity(targetTable.Name),
		strings.Join(columnsQuoted, ", "),
		strings.Join(columnsScopped, ", "),
		strings.Join(columnsUpdated, ", "),
	)
}

func (d *MySQLExecutorDialector) scopeNewColumn(column string) string {
	return fmt.Sprintf("NEW.%s", d.formatter.QuoteEntity(column))
}

func (d *MySQLExecutorDialector) scopeOldColumn(column string) string {
	return fmt.Sprintf("OLD.%s", d.formatter.QuoteEntity(column))
}

func (d *MySQLExecutorDialector) TableCreateQuery(sourceTable string, targetTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE TABLE %s LIKE %s`),
		d.formatter.ScopeEntity(d.schema, targetTable),
		d.formatter.ScopeEntity(d.schema, sourceTable),
	)
}

func (d *MySQLExecutorDialector) TableRemoveQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP TABLE IF EXISTS %s`),
		d.formatter.ScopeEntity(d.schema, sourceTable),
	)
}

func (d *MySQLExecutorDialector) TableInsertBatchQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string, lowerID int, upperID int) string {
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
		d.formatter.TrimSpaces(`INSERT IGNORE INTO %s (%s) %s`),
		scopedTargetTable,
		columnsJoined,
		selectQuery,
	)
}

func (d *MySQLExecutorDialector) TriggerCreateDeleteQuery(sourceTable models.TableInfo, targetTable models.TableInfo) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE TRIGGER %s AFTER DELETE ON %s FOR EACH ROW %s`),
		d.formatter.PrefixName(`_delete`, sourceTable.Name),
		d.formatter.ScopeEntity(d.schema, sourceTable.Name),
		d.buildDeleteQuery(targetTable),
	)
}

func (d *MySQLExecutorDialector) TriggerCreateInsertQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE TRIGGER %s AFTER INSERT ON %s FOR EACH ROW %s`),
		d.formatter.PrefixName(`_insert`, sourceTable.Name),
		d.formatter.ScopeEntity(d.schema, sourceTable.Name),
		d.buildUpsertQuery(targetTable, columns),
	)
}

func (d *MySQLExecutorDialector) TriggerCreateUpdateQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE TRIGGER %s AFTER UPDATE ON %s FOR EACH ROW %s`),
		d.formatter.PrefixName(`_update`, sourceTable.Name),
		d.formatter.ScopeEntity(d.schema, sourceTable.Name),
		d.buildUpsertQuery(targetTable, columns),
	)
}

func (d *MySQLExecutorDialector) TriggerDropDeleteQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP TRIGGER IF EXISTS %s`),
		d.formatter.PrefixName(`_delete`, sourceTable),
	)
}

func (d *MySQLExecutorDialector) TriggerDropInsertQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP TRIGGER IF EXISTS %s`),
		d.formatter.PrefixName(`_insert`, sourceTable),
	)
}

func (d *MySQLExecutorDialector) TriggerDropUpdateQuery(sourceTable string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP TRIGGER IF EXISTS %s`),
		d.formatter.PrefixName(`_update`, sourceTable),
	)
}

package migrator

import (
	"fmt"
	"journal-migrator/lib"
	"journal-migrator/pkg/storage/formatters"
	"strings"
)

type MySQLMigratorDialector struct {
	formatter *formatters.MySQLFormatter
	schema    string
}

func NewMySQLMigratorDialector(formatter *formatters.MySQLFormatter, schema string) *MySQLMigratorDialector {
	return &MySQLMigratorDialector{
		formatter: formatter,
		schema:    schema,
	}
}

func (d *MySQLMigratorDialector) buildColumnDefinition(dataType string, defaultVal any, nullable bool) string {
	var definitionParts = []string{dataType}

	if defaultVal != nil {
		columnDefault := d.formatter.QuoteValue(defaultVal)
		columnDefault = fmt.Sprintf("DEFAULT %s", columnDefault)
		definitionParts = append(definitionParts, columnDefault)
	}

	if nullable {
		definitionParts = append(definitionParts, "NULL")
	} else {
		definitionParts = append(definitionParts, "NOT NULL")
	}

	return strings.Join(definitionParts, " ")
}

func (d *MySQLMigratorDialector) CreateColumnQuery(table string, name string, dataType string, defaultVal any, nullable bool) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s ADD COLUMN %s %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(name),
		d.buildColumnDefinition(dataType, defaultVal, nullable),
	)
}

func (d *MySQLMigratorDialector) ChangeColumnQuery(table string, name string, dataType string, defaultVal any, nullable bool) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s MODIFY COLUMN %s %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(name),
		d.buildColumnDefinition(dataType, defaultVal, nullable),
	)
}

func (d *MySQLMigratorDialector) RemoveColumnQuery(table string, name string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s DROP COLUMN %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(name),
	)
}

func (d *MySQLMigratorDialector) RenameColumnQuery(table string, oldName string, newName string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s RENAME COLUMN %s TO %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(oldName),
		d.formatter.QuoteEntity(newName),
	)
}

func (d *MySQLMigratorDialector) CreateIndexQuery(table string, name string, columns []string, unique bool) string {
	var query string

	if unique {
		query = `
			ALTER TABLE %s
			ADD UNIQUE INDEX %s (%s)
		`
	} else {
		query = `
			ALTER TABLE %s
			ADD INDEX %s (%s)
		`
	}

	var indexColumnsQuoted = lib.Mapper(columns, d.formatter.QuoteEntity)
	var indexColumnsJoined = strings.Join(indexColumnsQuoted, ", ")

	return fmt.Sprintf(
		d.formatter.TrimSpaces(query),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(name),
		indexColumnsJoined,
	)
}

func (d *MySQLMigratorDialector) ChangeIndexQuery(_ string, _ string, _ []string, _ bool) string {
	panic("Operation not supported")
}

func (d *MySQLMigratorDialector) RemoveIndexQuery(table string, name string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s DROP INDEX %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(name),
	)
}

func (d *MySQLMigratorDialector) RenameIndexQuery(table string, oldName string, newName string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s RENAME INDEX %s TO %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(oldName),
		d.formatter.QuoteEntity(newName),
	)
}

func (d *MySQLMigratorDialector) CreateTableQuery(name string, definition string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE TABLE %s (%s)`),
		d.formatter.ScopeEntity(d.schema, name),
		definition,
	)
}

func (d *MySQLMigratorDialector) RemoveTableQuery(name string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP TABLE %s`),
		d.formatter.ScopeEntity(d.schema, name),
	)
}

func (d *MySQLMigratorDialector) RenameTableQuery(oldName, newName string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`RENAME TABLE %s.%s TO %s`),
		d.formatter.QuoteEntity(d.schema),
		d.formatter.QuoteEntity(oldName),
		d.formatter.QuoteEntity(newName),
	)
}

func (d *MySQLMigratorDialector) DDLOperationQuery(statement string) string {
	return statement
}

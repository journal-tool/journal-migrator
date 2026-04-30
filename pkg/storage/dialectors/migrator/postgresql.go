package migrator

import (
	"fmt"
	"journal-migrator/lib"
	"journal-migrator/pkg/storage/formatters"
	"strings"
)

type PostgreSQLMigratorDialector struct {
	formatter *formatters.PostgreSQLFormatter
	schema    string
}

func NewPostgreSQLMigratorDialector(formatter *formatters.PostgreSQLFormatter, schema string) *PostgreSQLMigratorDialector {
	return &PostgreSQLMigratorDialector{
		formatter: formatter,
		schema:    schema,
	}
}

func (d *PostgreSQLMigratorDialector) buildCreateColumnDefinition(dataType string, defaultVal any, nullable bool) string {
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

func (d *PostgreSQLMigratorDialector) buildChangeColumnModifier(name string, dataType string, defaultVal any, nullable bool) string {
	var modifiers []string
	var query = fmt.Sprintf(`ALTER COLUMN %s %%s`, d.formatter.QuoteEntity(name))

	if dataType != "" {
		modifier := fmt.Sprintf(`SET DATA TYPE %s USING (%s::%s)`, dataType, name, dataType)
		modifier = fmt.Sprintf(query, modifier)
		modifiers = append(modifiers, modifier)
	}

	if defaultVal != nil {
		modifier := fmt.Sprintf(`SET DEFAULT %s`, d.formatter.QuoteValue(defaultVal))
		modifier = fmt.Sprintf(query, modifier)
		modifiers = append(modifiers, modifier)
	} else {
		modifier := fmt.Sprintf(`DROP DEFAULT`)
		modifier = fmt.Sprintf(query, modifier)
		modifiers = append(modifiers, modifier)
	}

	if nullable {
		modifier := fmt.Sprintf(`SET NOT NULL`)
		modifier = fmt.Sprintf(query, modifier)
		modifiers = append(modifiers, modifier)
	} else {
		modifier := fmt.Sprintf(`DROP NOT NULL`)
		modifier = fmt.Sprintf(query, modifier)
		modifiers = append(modifiers, modifier)
	}

	return strings.Join(modifiers, ", ")
}

func (d *PostgreSQLMigratorDialector) CreateColumnQuery(table string, name string, dataType string, defaultVal any, nullable bool) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s ADD COLUMN %s %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(name),
		d.buildCreateColumnDefinition(dataType, defaultVal, nullable),
	)
}

func (d *PostgreSQLMigratorDialector) ChangeColumnQuery(table string, name string, dataType string, defaultVal any, nullable bool) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.buildChangeColumnModifier(name, dataType, defaultVal, nullable),
	)
}

func (d *PostgreSQLMigratorDialector) RemoveColumnQuery(table string, name string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s DROP COLUMN %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(name),
	)
}

func (d *PostgreSQLMigratorDialector) RenameColumnQuery(table string, oldName string, newName string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s RENAME COLUMN %s TO %s`),
		d.formatter.ScopeEntity(d.schema, table),
		d.formatter.QuoteEntity(oldName),
		d.formatter.QuoteEntity(newName),
	)
}

func (d *PostgreSQLMigratorDialector) CreateIndexQuery(table string, name string, columns []string, unique bool) string {
	var query string

	if unique {
		query = `
			CREATE UNIQUE INDEX %s
			ON %s (%s)
		`
	} else {
		query = `
			CREATE INDEX %s
			ON %s (%s)
		`
	}

	var indexColumnsQuoted = lib.Mapper(columns, d.formatter.QuoteEntity)
	var indexColumnsJoined = strings.Join(indexColumnsQuoted, ", ")

	return fmt.Sprintf(
		d.formatter.TrimSpaces(query),
		d.formatter.QuoteEntity(name),
		d.formatter.ScopeEntity(d.schema, table),
		indexColumnsJoined,
	)
}

func (d *PostgreSQLMigratorDialector) ChangeIndexQuery(_ string, _ string, _ []string, _ bool) string {
	panic("Operation not supported")
}

func (d *PostgreSQLMigratorDialector) RemoveIndexQuery(_ string, name string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP INDEX %s`),
		d.formatter.QuoteEntity(name),
	)
}

func (d *PostgreSQLMigratorDialector) RenameIndexQuery(_ string, oldName string, newName string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER INDEX %s RENAME TO %s`),
		d.formatter.QuoteEntity(oldName),
		d.formatter.QuoteEntity(newName),
	)
}

func (d *PostgreSQLMigratorDialector) CreateTableQuery(name string, definition string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`CREATE TABLE %s (%s)`),
		d.formatter.ScopeEntity(d.schema, name),
		definition,
	)
}

func (d *PostgreSQLMigratorDialector) RemoveTableQuery(name string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`DROP TABLE %s`),
		d.formatter.ScopeEntity(d.schema, name),
	)
}

func (d *PostgreSQLMigratorDialector) RenameTableQuery(oldName string, newName string) string {
	return fmt.Sprintf(
		d.formatter.TrimSpaces(`ALTER TABLE %s.%s RENAME TO %s`),
		d.formatter.QuoteEntity(d.schema),
		d.formatter.QuoteEntity(oldName),
		d.formatter.QuoteEntity(newName),
	)
}

func (d *PostgreSQLMigratorDialector) DDLOperationQuery(statement string) string {
	return statement
}

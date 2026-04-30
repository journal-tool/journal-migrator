package fetcher

import (
	"fmt"
	"journal-migrator/pkg/storage/formatters"
)

type MySQLFetcherDialector struct {
	formatter *formatters.MySQLFormatter
	schema    string
}

func NewMySQLFetcherDialector(formatter *formatters.MySQLFormatter, schema string) *MySQLFetcherDialector {
	return &MySQLFetcherDialector{
		formatter: formatter,
		schema:    schema,
	}
}

func (d *MySQLFetcherDialector) ColumnSelectMaxQuery(table string, column string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			COALESCE(MAX(%s), 0)
		FROM
			%s
	`)

	return fmt.Sprintf(query, column, d.formatter.ScopeEntity(d.schema, table))
}

func (d *MySQLFetcherDialector) TableSelectInfoQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			table_name,
			column_name,
			extra = 'auto_increment'
		FROM
			information_schema.columns
		JOIN
			information_schema.key_column_usage USING (table_schema, table_name, column_name)
		WHERE
			constraint_name = 'PRIMARY' AND
			table_schema = '%s' AND
			table_name = '%s'
	`)

	return fmt.Sprintf(query, d.schema, table)
}

func (d *MySQLFetcherDialector) TableSelectSizeQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			table_name,
			data_length,
			index_length
		FROM
			information_schema.tables
		WHERE
			table_schema = '%s' AND
			table_name = '%s'
	`)

	return fmt.Sprintf(query, d.schema, table)
}

func (d *MySQLFetcherDialector) TableSelectColumnsQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			column_name,
			column_type,
			column_default,
			collation_name,
			is_nullable = 'YES'
		FROM
			information_schema.columns
		WHERE
			table_schema = '%s' AND
			table_name = '%s'
		ORDER BY
			ORDINAL_POSITION
	`)

	return fmt.Sprintf(query, d.schema, table)
}

func (d *MySQLFetcherDialector) TableSelectIndexesQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			index_name,
			index_type,
			non_unique = 0,
			GROUP_CONCAT(column_name SEPARATOR ', ')
		FROM
			information_schema.statistics
		WHERE
			table_schema = '%s' AND
			table_name = '%s'
		GROUP BY
			index_name,
			index_type,
			non_unique
	`)

	return fmt.Sprintf(query, d.schema, table)
}

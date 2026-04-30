package fetcher

import (
	"fmt"
	"journal-migrator/pkg/storage/formatters"
)

type PostgreSQLFetcherDialector struct {
	formatter *formatters.PostgreSQLFormatter
	schema    string
}

func NewPostgreSQLFetcherDialector(formatter *formatters.PostgreSQLFormatter, schema string) *PostgreSQLFetcherDialector {
	return &PostgreSQLFetcherDialector{
		formatter: formatter,
		schema:    schema,
	}
}

func (d *PostgreSQLFetcherDialector) ColumnSelectMaxQuery(table string, column string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			COALESCE(MAX(%s), 0)
		FROM
			%s
	`)

	return fmt.Sprintf(query, column, d.formatter.ScopeEntity(d.schema, table))
}

func (d *PostgreSQLFetcherDialector) TableSelectInfoQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			table_name,
			column_name,
			is_identity = 'YES'
		FROM
			information_schema.columns
		JOIN
			information_schema.table_constraints USING (table_schema, table_name)
		JOIN
			information_schema.constraint_column_usage USING (table_schema, table_name, constraint_name, column_name)
		WHERE
			constraint_type = 'PRIMARY KEY' AND
			table_schema = '%s' AND
			table_name = '%s'
	`)

	return fmt.Sprintf(query, d.schema, table)
}

func (d *PostgreSQLFetcherDialector) TableSelectSizeQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			table_name,
			pg_table_size(table_schema || '.' || table_name),
			pg_indexes_size(table_schema || '.' || table_name)
		FROM
			information_schema.tables
		WHERE
			table_type = 'BASE TABLE' AND
			table_schema = '%s' AND
			table_name = '%s'
	`)

	return fmt.Sprintf(query, d.schema, table)
}

func (d *PostgreSQLFetcherDialector) TableSelectColumnsQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			column_name,
			data_type,
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

func (d *PostgreSQLFetcherDialector) TableSelectIndexesQuery(table string) string {
	var query = d.formatter.TrimSpaces(`
		SELECT
			indexname,
			(regexp_match(indexdef, 'USING (\w+)', 'i'))[1],
			(regexp_like(indexdef, 'UNIQUE INDEX', 'i')),
			(regexp_match(indexdef, '\((.*)\)', 'i'))[1]
		FROM
			pg_catalog.pg_indexes
		WHERE
			schemaname = '%s' AND
			tablename = '%s'
	`)

	return fmt.Sprintf(query, d.schema, table)
}

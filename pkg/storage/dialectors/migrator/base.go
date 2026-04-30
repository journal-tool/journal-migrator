package migrator

type BaseMigratorDialector interface {
	CreateColumnQuery(table string, name string, dataType string, defaultVal any, nullable bool) string
	ChangeColumnQuery(table string, name string, dataType string, defaultVal any, nullable bool) string
	RemoveColumnQuery(table string, name string) string
	RenameColumnQuery(table string, oldName string, newName string) string

	CreateIndexQuery(table string, name string, columns []string, unique bool) string
	ChangeIndexQuery(table string, name string, columns []string, unique bool) string
	RemoveIndexQuery(table string, name string) string
	RenameIndexQuery(table string, oldName string, newName string) string

	CreateTableQuery(name string, definition string) string
	RemoveTableQuery(name string) string
	RenameTableQuery(oldName string, newName string) string

	DDLOperationQuery(statement string) string
}

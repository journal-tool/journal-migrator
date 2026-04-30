package operation

const (
	CreateTableType = "CREATE_TABLE"
	RemoveTableType = "REMOVE_TABLE"
	RenameTableType = "RENAME_TABLE"

	CreateColumnType = "CREATE_COLUMN"
	ChangeColumnType = "CHANGE_COLUMN"
	RemoveColumnType = "REMOVE_COLUMN"
	RenameColumnType = "RENAME_COLUMN"

	CreateIndexType = "CREATE_INDEX"
	ChangeIndexType = "CHANGE_INDEX"
	RemoveIndexType = "REMOVE_INDEX"
	RenameIndexType = "RENAME_INDEX"

	DDLType = "DDL"
)

var TableOperations = []string{
	CreateTableType,
	RemoveTableType,
	RenameTableType,
}

package schema

import "database/sql"

type ColumnInfo struct {
	Name      string
	Type      string
	Default   sql.NullString
	Collation sql.NullString
	Nullable  bool
}

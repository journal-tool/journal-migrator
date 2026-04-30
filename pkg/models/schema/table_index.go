package schema

import "strings"

type IndexInfo struct {
	Name    string
	Type    string
	Unique  bool
	Columns string
}

func (i IndexInfo) ColumnNames() []string {
	return strings.Split(i.Columns, ", ")
}

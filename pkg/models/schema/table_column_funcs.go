package schema

import "journal-migrator/lib"

func extractColumnName(column ColumnInfo) string {
	return column.Name
}

func IntersectColumns(columnInfoList ...[]ColumnInfo) []string {
	var columnNameLists [][]string

	for _, columnInfos := range columnInfoList {
		columnNameList := lib.Mapper(columnInfos, extractColumnName)
		columnNameLists = append(columnNameLists, columnNameList)
	}

	return lib.Intersect(columnNameLists...)
}

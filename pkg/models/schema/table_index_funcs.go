package schema

import "journal-migrator/lib"

func extractIndexName(Index IndexInfo) string {
	return Index.Name
}

func IntersectIndexes(indexInfoList ...[]IndexInfo) []string {
	var indexNameLists [][]string

	for _, indexInfos := range indexInfoList {
		indexNameList := lib.Mapper(indexInfos, extractIndexName)
		indexNameLists = append(indexNameLists, indexNameList)
	}

	return lib.Intersect(indexNameLists...)
}

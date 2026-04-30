package fetcher

type BaseFetcherDialector interface {
	ColumnSelectMaxQuery(table string, column string) string

	TableSelectInfoQuery(table string) string
	TableSelectSizeQuery(table string) string
	TableSelectColumnsQuery(table string) string
	TableSelectIndexesQuery(table string) string
}

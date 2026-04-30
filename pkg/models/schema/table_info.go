package schema

type TableInfo struct {
	Name  string
	Key   string
	Valid bool
}

type TableSize struct {
	Name      string
	DataSize  int64
	IndexSize int64
}

func (t TableSize) TotalSize() int64 {
	return t.DataSize + t.IndexSize
}

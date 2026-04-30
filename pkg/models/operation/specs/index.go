package specs

type CreateIndexSpec struct {
	Name    string   `json:"name"`
	Unique  bool     `json:"unique"`
	Columns []string `json:"columns"`
}

type ChangeIndexSpec struct {
	Name    string   `json:"name"`
	Unique  bool     `json:"unique"`
	Columns []string `json:"columns"`
}

type RenameIndexSpec struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type RemoveIndexSpec struct {
	Name string `json:"name"`
}

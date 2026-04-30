package specs

type CreateTableSpec struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type RenameTableSpec struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type RemoveTableSpec struct {
	Name string `json:"name"`
}

package specs

type CreateColumnSpec struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Default  any    `json:"default"`
	Nullable bool   `json:"nullable"`
}

type ChangeColumnSpec struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Default  any    `json:"default"`
	Nullable bool   `json:"nullable"`
}

type RenameColumnSpec struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type RemoveColumnSpec struct {
	Name string `json:"name"`
}

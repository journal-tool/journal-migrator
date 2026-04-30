package operation

import (
	"bytes"
	"encoding/json"
	"journal-migrator/pkg/models/operation/specs"
	"slices"
)

type RawOperation struct {
	Type string          `json:"type"`
	Spec json.RawMessage `json:"spec"`
}

type Operation struct {
	Type string `json:"type"`
	Spec any    `json:"spec"`
}

func (op *Operation) IsTableOperation() bool {
	return slices.Contains(TableOperations, op.Type)
}

func (op *Operation) UnmarshalJSON(data []byte) error {
	var raw RawOperation

	reader := bytes.NewReader(data)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	var err = decoder.Decode(&raw)
	if err != nil {
		return err
	}

	var spec any
	switch raw.Type {
	case CreateColumnType:
		spec = &specs.CreateColumnSpec{}
	case ChangeColumnType:
		spec = &specs.ChangeColumnSpec{}
	case RemoveColumnType:
		spec = &specs.RemoveColumnSpec{}
	case RenameColumnType:
		spec = &specs.RenameColumnSpec{}
	case CreateIndexType:
		spec = &specs.CreateIndexSpec{}
	case ChangeIndexType:
		spec = &specs.ChangeIndexSpec{}
	case RemoveIndexType:
		spec = &specs.RemoveIndexSpec{}
	case RenameIndexType:
		spec = &specs.RenameIndexSpec{}
	case CreateTableType:
		spec = &specs.CreateTableSpec{}
	case RemoveTableType:
		spec = &specs.RemoveTableSpec{}
	case RenameTableType:
		spec = &specs.RenameTableSpec{}
	case DDLType:
		spec = &specs.DDLSpec{}
	}

	reader = bytes.NewReader(raw.Spec)
	decoder = json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(spec)
	if err != nil {
		return err
	}

	op.Type = raw.Type
	op.Spec = spec
	return nil
}

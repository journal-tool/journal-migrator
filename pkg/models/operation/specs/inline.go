package specs

import (
	"encoding/json"
	"errors"
	"strings"
)

type DDLSpec struct {
	Statement string `json:"statement"`
}

func (s *DDLSpec) UnmarshalJSON(data []byte) error {
	type tmpSpec DDLSpec

	var spec tmpSpec
	var err = json.Unmarshal(data, &spec)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(spec.Statement, "ALTER TABLE") {
		return errors.New("statement must be an ALTER TABLE one")
	}

	s.Statement = spec.Statement
	return nil
}

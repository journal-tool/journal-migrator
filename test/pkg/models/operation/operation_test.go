package operation

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/models/operation"
	"journal-migrator/pkg/models/operation/specs"
	"testing"
)

func TestOperation(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping unit tests")
	}

	suite.Run(t, new(OperationSuite))
}

type OperationSuite struct {
	suite.Suite
}

func (suite *OperationSuite) assertPanic(testFn func()) {
	defer func() {
		r := recover()
		suite.NotNil(r)
	}()

	testFn()
}

func (suite *OperationSuite) TestParseValidOperations() {
	raw := `[
		{"type": "CREATE_TABLE", "spec": {"name": "test_table", "definition": "id SERIAL PRIMARY KEY"}},
		{"type": "RENAME_COLUMN", "spec": {"old_name": "old_column", "new_name": "new_name"}}
	]`

	ops := models.ParseOperations(raw)

	suite.Equal(operation.CreateTableType, ops[0].Type)
	suite.IsType(&specs.CreateTableSpec{}, ops[0].Spec)

	suite.Equal(operation.RenameColumnType, ops[1].Type)
	suite.IsType(&specs.RenameColumnSpec{}, ops[1].Spec)
}

func (suite *OperationSuite) TestParseInvalidOperations() {
	suite.assertPanic(func() {
		models.ParseOperations(`[{"unkown_field": "unknown_value"}]`)
	})

	suite.assertPanic(func() {
		models.ParseOperations(`[{"type": "INVALID_OP", "spec": {}}]`)
	})

	suite.assertPanic(func() {
		models.ParseOperations(`[{"type": "CREATE_TABLE", "spec": {"unknown_field": "unknown_value"}}]`)
	})
}

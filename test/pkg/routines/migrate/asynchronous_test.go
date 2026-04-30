package migrate

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/models/operation"
	"journal-migrator/pkg/models/operation/specs"
	"journal-migrator/pkg/routines"
	"journal-migrator/pkg/storage"
	"journal-migrator/test"
	"testing"
)

func TestAsyncMigrateRoutineSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(AsyncMigrateRoutineSuite))
}

type AsyncMigrateRoutineSuite struct {
	test.IntegrationSuite
}

func (suite *AsyncMigrateRoutineSuite) createRoutine() routines.BaseMigrateRoutine {
	return routines.NewMigrateRoutine(suite.Client, suite.Logger, routines.AsyncStrategy)
}

func (suite *AsyncMigrateRoutineSuite) TestRunInvalidOperations() {
	tableName := "test_table"
	tableOps := []models.Operation{
		{Type: operation.CreateTableType, Spec: &specs.CreateTableSpec{}},
	}

	routine := suite.createRoutine()

	err := routine.Run(tableName, tableOps, nil)
	suite.Error(err)
	suite.ErrorContains(err, "cannot perform table operations asynchronously")
}

func (suite *AsyncMigrateRoutineSuite) TestRunInvalidTableKey() {
	tableName := "test_table"
	tableDef := "(id BIGINT PRIMARY KEY, data TEXT)"
	tableOps := []models.Operation{
		{Type: operation.CreateColumnType, Spec: &specs.CreateColumnSpec{}},
	}

	routine := suite.createRoutine()
	database := suite.Client.GetDatabase()

	suite.WithTable(tableName, tableDef, database, func(db storage.Database) {
		err := routine.Run(tableName, tableOps, nil)
		suite.Error(err)
		suite.ErrorContains(err, "table does not have a valid auto-incremented key")
	})
}

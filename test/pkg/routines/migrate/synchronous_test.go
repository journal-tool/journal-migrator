package migrate

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/models/operation"
	"journal-migrator/pkg/models/operation/specs"
	"journal-migrator/pkg/routines"
	"journal-migrator/test"
	"testing"
)

func TestSyncMigrateRoutineSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(SyncMigrateRoutineSuite))
}

type SyncMigrateRoutineSuite struct {
	test.IntegrationSuite
}

func (suite *SyncMigrateRoutineSuite) createRoutine() routines.BaseMigrateRoutine {
	return routines.NewMigrateRoutine(suite.Client, suite.Logger, routines.SyncStrategy)
}

func (suite *SyncMigrateRoutineSuite) TestRunInvalidOperations() {
	tableName := "test_table"
	tableOps := []models.Operation{
		{Type: operation.CreateTableType, Spec: &specs.CreateTableSpec{}},
		{Type: operation.RemoveTableType, Spec: &specs.RemoveTableSpec{}},
	}

	routine := suite.createRoutine()

	err := routine.Run(tableName, tableOps, nil)
	suite.Error(err)
	suite.ErrorContains(err, "cannot perform table operations in bulk")
}

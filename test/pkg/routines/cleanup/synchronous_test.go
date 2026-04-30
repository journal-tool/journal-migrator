package cleanup

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/routines"
	"journal-migrator/pkg/routines/routines/cleanup"
	"journal-migrator/pkg/storage"
	"journal-migrator/test"
	"testing"
)

func TestSyncCleanupRoutineSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(SyncCleanupRoutineSuite))
}

type SyncCleanupRoutineSuite struct {
	test.IntegrationSuite
}

func (suite *SyncCleanupRoutineSuite) createRoutine() cleanup.BaseCleanupRoutine {
	return routines.NewCleanupRoutine(
		suite.Client,
		suite.Logger,
		routines.SyncStrategy,
	)
}

func (suite *SyncCleanupRoutineSuite) TestRunWithoutTables() {
	routine := suite.createRoutine()

	err := routine.Run("table")
	suite.Nil(err)
}

func (suite *SyncCleanupRoutineSuite) TestRunWithTables() {
	tableName := "_migrator_table"
	tableDef := suite.GetTableDefinition()

	routine := suite.createRoutine()
	database := suite.Client.GetDatabase()

	suite.WithTable(tableName, tableDef, database, func(db storage.Database) {
		err := routine.Run("table")
		suite.Nil(err)
	})
}

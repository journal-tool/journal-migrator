package dialectors

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"journal-migrator/test"
	"testing"
)

func TestExecutorDialectorSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(ExecutorDialectorSuite))
}

type ExecutorDialectorSuite struct {
	test.IntegrationSuite
}

func (suite *ExecutorDialectorSuite) createDialector() dialectors.BaseExecutorDialector {
	config := suite.Client.GetConfig()

	formatter := formatters.NewFormatter(config.Dialect)
	dialector := dialectors.NewExecutorDialector(formatter, config.Name)

	return dialector
}

func (suite *ExecutorDialectorSuite) TestCreateTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	sourceTable := "test_source"
	targetTable := "test_target"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()
	createTableQuery := dialector.TableCreateQuery(sourceTable, targetTable)
	deleteTableQuery := dialector.TableRemoveQuery(targetTable)

	suite.WithTable(sourceTable, tableDef, tx, func(db storage.Database) {
		_, err := tx.Exec(createTableQuery)
		suite.Nil(err)
		_, err = tx.Exec(deleteTableQuery)
		suite.Nil(err)
	})
}

func (suite *ExecutorDialectorSuite) TestTableInsertBatchQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	dialector := suite.createDialector()
	insertQuery := dialector.TableInsertBatchQuery(sourceTable, targetTable, tableCols, 0000, 1000)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			_, err := tx.Exec(insertQuery)
			suite.Nil(err)
		})
	})
}

func (suite *ExecutorDialectorSuite) TestTriggerCreateDeleteQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	dialector := suite.createDialector()
	triggerQuery := dialector.TriggerCreateDeleteQuery(sourceTable, targetTable)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			_, err := tx.Exec(triggerQuery)
			suite.Nil(err)
		})
	})
}

func (suite *ExecutorDialectorSuite) TestTriggerCreateInsertQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	dialector := suite.createDialector()
	triggerQuery := dialector.TriggerCreateInsertQuery(sourceTable, targetTable, tableCols)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			_, err := tx.Exec(triggerQuery)
			suite.Nil(err)
		})
	})
}

func (suite *ExecutorDialectorSuite) TestTriggerCreateUpdateQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	dialector := suite.createDialector()
	triggerQuery := dialector.TriggerCreateUpdateQuery(sourceTable, targetTable, tableCols)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			_, err := tx.Exec(triggerQuery)
			suite.Nil(err)
		})
	})
}

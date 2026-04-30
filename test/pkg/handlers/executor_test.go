package handlers

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/handlers"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"journal-migrator/test"
	"testing"
)

func TestExecutorHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(ExecutorHandlerSuite))
}

type ExecutorHandlerSuite struct {
	test.IntegrationSuite
}

func (suite *ExecutorHandlerSuite) createHandler(database storage.Database) *handlers.ExecutorHandler {
	config := *suite.Client.GetConfig()

	formatter := formatters.NewFormatter(config.Dialect)
	dialector := dialectors.NewExecutorDialector(formatter, config.Name)
	handler := handlers.NewExecutorHandler(dialector, database, config, suite.Logger)

	return handler
}

func (suite *ExecutorHandlerSuite) validateRows(tx *sql.Tx, sourceTable string, targetTable string) {
	parseFunc := suite.GetTableParseFunc()
	sourceTableQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY id", sourceTable)
	targetTableQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY id", targetTable)

	sourceResult, err := tx.Query(sourceTableQuery)
	sourceRows, err := parseFunc(sourceResult)
	suite.Nil(err)

	targetResult, err := tx.Query(targetTableQuery)
	targetRows, err := parseFunc(targetResult)
	suite.Nil(err)

	suite.Equal(sourceRows, targetRows)
}

func (suite *ExecutorHandlerSuite) TestInsertRowsEmptyTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	handler := suite.createHandler(tx)

	sourceInsertValue := "(1, 'text-default')"
	sourceInsertQuery := fmt.Sprintf("INSERT INTO %s VALUES %s", sourceTable.Name, sourceInsertValue)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			_, _ = tx.Exec(sourceInsertQuery)

			err := handler.InsertRows(sourceTable, targetTable, tableCols, 1, 10)
			suite.Nil(err)

			suite.validateRows(tx, sourceTable.Name, targetTable.Name)
		})
	})
}

func (suite *ExecutorHandlerSuite) TestInsertRowsFilledTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	handler := suite.createHandler(tx)

	sourceInsertValue := "(1, 'text-outdated')"
	targetInsertValue := "(1, 'text-replicated')"

	sourceInsertQuery := fmt.Sprintf("INSERT INTO %s VALUES %s", sourceTable.Name, sourceInsertValue)
	targetInsertQuery := fmt.Sprintf("INSERT INTO %s VALUES %s", targetTable.Name, targetInsertValue)
	targetSelectQuery := fmt.Sprintf("SELECT data FROM %s WHERE id = 1", targetTable.Name)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			_, _ = tx.Exec(sourceInsertQuery)
			_, _ = tx.Exec(targetInsertQuery)

			err := handler.InsertRows(sourceTable, targetTable, tableCols, 1, 10)
			suite.Nil(err)

			data := ""
			row := tx.QueryRow(targetSelectQuery)
			err = row.Scan(&data)
			suite.Nil(err)
			suite.Equal("text-replicated", data)
		})
	})
}

func (suite *ExecutorHandlerSuite) TestCreateTriggersEmptyTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	handler := suite.createHandler(tx)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			var query string
			var err error

			err = handler.CreateTableTriggers(sourceTable, targetTable, tableCols)
			suite.Nil(err)

			query = fmt.Sprintf("INSERT INTO %s VALUES (1, 'text-inserted')", sourceTable.Name)
			_, err = tx.Exec(query)
			suite.Nil(err)
			suite.validateRows(tx, sourceTable.Name, targetTable.Name)

			query = fmt.Sprintf("UPDATE %s SET data = 'text-updated' WHERE id = 1", sourceTable.Name)
			_, err = tx.Exec(query)
			suite.Nil(err)
			suite.validateRows(tx, sourceTable.Name, targetTable.Name)

			query = fmt.Sprintf("DELETE FROM %s WHERE id = 1", sourceTable.Name)
			_, err = tx.Exec(query)
			suite.Nil(err)
			suite.validateRows(tx, sourceTable.Name, targetTable.Name)
		})
	})
}

func (suite *ExecutorHandlerSuite) TestCreateTriggersFilledTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	handler := suite.createHandler(tx)

	targetInsertValue := "(1, 'text-outdated')"
	targetInsertQuery := fmt.Sprintf("INSERT INTO %s VALUES %s", targetTable.Name, targetInsertValue)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			var query string
			var err error

			_, _ = tx.Exec(targetInsertQuery)

			err = handler.CreateTableTriggers(sourceTable, targetTable, tableCols)
			suite.Nil(err)

			query = fmt.Sprintf("INSERT INTO %s VALUES (1, 'text-inserted')", sourceTable.Name)
			_, err = tx.Exec(query)
			suite.Nil(err)
			suite.validateRows(tx, sourceTable.Name, targetTable.Name)

			query = fmt.Sprintf("UPDATE %s SET data = 'text-updated' WHERE id = 1", sourceTable.Name)
			_, err = tx.Exec(query)
			suite.Nil(err)
			suite.validateRows(tx, sourceTable.Name, targetTable.Name)

			query = fmt.Sprintf("DELETE FROM %s WHERE id = 1", sourceTable.Name)
			_, err = tx.Exec(query)
			suite.Nil(err)
			suite.validateRows(tx, sourceTable.Name, targetTable.Name)
		})
	})
}

func (suite *ExecutorHandlerSuite) TestDeleteTriggers() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableCols := suite.GetTableColumnNames()

	sourceTable := models.TableInfo{Name: "test_source", Key: tableKey}
	targetTable := models.TableInfo{Name: "test_target", Key: tableKey}

	handler := suite.createHandler(tx)

	suite.WithTable(sourceTable.Name, tableDef, tx, func(db storage.Database) {
		suite.WithTable(targetTable.Name, tableDef, tx, func(db storage.Database) {
			err := handler.DeleteTableTriggers(sourceTable.Name)
			suite.Nil(err)

			err = handler.CreateTableTriggers(sourceTable, targetTable, tableCols)
			err = handler.DeleteTableTriggers(sourceTable.Name)
			suite.Nil(err)
		})
	})
}

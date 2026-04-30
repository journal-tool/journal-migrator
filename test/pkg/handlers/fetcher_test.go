package handlers

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/handlers"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"journal-migrator/test"
	"slices"
	"strings"
	"testing"
)

func TestFetcherHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(FetcherHandlerSuite))
}

type FetcherHandlerSuite struct {
	test.IntegrationSuite
}

func (suite *FetcherHandlerSuite) createHandler(database storage.Database) *handlers.FetcherHandler {
	config := *suite.Client.GetConfig()

	formatter := formatters.NewFormatter(config.Dialect)
	dialector := dialectors.NewFetcherDialector(formatter, config.Name)
	handler := handlers.NewFetcherHandler(dialector, database, config, suite.Logger)

	return handler
}

func (suite *FetcherHandlerSuite) TestFetchColumnMaxFilledTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()
	tableLen := 10

	handler := suite.createHandler(tx)

	insertQuery := "INSERT INTO test_table (data) VALUES %s"
	insertValues := strings.Join(slices.Repeat([]string{"('text')"}, tableLen), ", ")

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Exec(fmt.Sprintf(insertQuery, insertValues))
		suite.Nil(err)

		maxKey, err := handler.FetchColumnMax(tableName, tableKey)
		suite.Nil(err)
		suite.Equal(tableLen, maxKey)
	})
}

func (suite *FetcherHandlerSuite) TestFetchColumnMaxEmptyTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()

	handler := suite.createHandler(tx)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		maxKey, err := handler.FetchColumnMax(tableName, tableKey)
		suite.Nil(err)
		suite.Equal(0, maxKey)
	})
}

func (suite *FetcherHandlerSuite) TestFetchTableInfoPresentTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()

	handler := suite.createHandler(tx)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		info, err := handler.FetchTableInfo(tableName)
		suite.Nil(err)

		suite.Equal(tableName, info.Name)
		suite.Equal(tableKey, info.Key)
		suite.True(info.Valid)
	})
}

func (suite *FetcherHandlerSuite) TestFetchTableInfoMissingTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	handler := suite.createHandler(tx)

	info, err := handler.FetchTableInfo(tableName)
	suite.NotNil(err)
	suite.Nil(info)
}

func (suite *FetcherHandlerSuite) TestFetchTableSize() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	handler := suite.createHandler(tx)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		size, err := handler.FetchTableSize(tableName)
		suite.Nil(err)

		suite.Equal(tableName, size.Name)
		suite.NotZero(size.DataSize)
	})
}

func (suite *FetcherHandlerSuite) TestFetchTableColumns() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	handler := suite.createHandler(tx)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		columns, err := handler.FetchTableColumns(tableName)
		suite.Nil(err)
		suite.Equal(2, len(columns))

		suite.Equal("id", columns[0].Name)
		suite.Equal("data", columns[1].Name)
	})
}

func (suite *FetcherHandlerSuite) TestFetchTableIndexes() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	handler := suite.createHandler(tx)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		indexes, err := handler.FetchTableIndexes(tableName)
		suite.Nil(err)
		suite.NotZero(len(indexes))
	})
}

func (suite *FetcherHandlerSuite) TestFetchTableColumnIntersection() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	sourceTableName := "test_source"
	targetTableName := "test_target"
	tableDefinition := suite.GetTableDefinition()
	tableColumns := suite.GetTableColumnNames()

	handler := suite.createHandler(tx)

	suite.WithTable(sourceTableName, tableDefinition, tx, func(db storage.Database) {
		suite.WithTable(targetTableName, tableDefinition, tx, func(db storage.Database) {
			sharedColumns, err := handler.FetchTableColumnIntersection(sourceTableName, targetTableName)
			suite.Nil(err)
			suite.Equal(tableColumns, sharedColumns)
		})
	})
}

//func (suite *FetcherHandlerSuite) TestFetchTableIndexIntersection() {
//	db := suite.Client.GetDatabase()
//	tx, _ := db.Begin()
//	defer tx.Rollback()
//
//	sourceTableName := "test_source"
//	targetTableName := "test_target"
//	tableDefinition := suite.GetTableDefinition()
//
//	handler := suite.createHandler(tx)
//
//	suite.WithTable(sourceTableName, tableDefinition, tx, func(db storage.Database) {
//		suite.WithTable(targetTableName, tableDefinition, tx, func(db storage.Database) {
//			sharedIndexes, err := handler.FetchTableIndexIntersection(sourceTableName, targetTableName)
//			suite.Nil(err)
//			suite.NotEmpty(sharedIndexes)
//		})
//	})
//}

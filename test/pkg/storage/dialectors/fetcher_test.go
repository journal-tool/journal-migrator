package dialectors

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"journal-migrator/test"
	"testing"
)

func TestFetcherDialectorSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(FetcherDialectorSuite))
}

type FetcherDialectorSuite struct {
	test.IntegrationSuite
}

func (suite *FetcherDialectorSuite) createDialector() dialectors.BaseFetcherDialector {
	config := suite.Client.GetConfig()

	formatter := formatters.NewFormatter(config.Dialect)
	dialector := dialectors.NewFetcherDialector(formatter, config.Name)

	return dialector
}

func (suite *FetcherDialectorSuite) TestSelectMaxQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()
	tableKey := suite.GetTablePrimaryKey()

	dialector := suite.createDialector()
	selectQuery := dialector.ColumnSelectMaxQuery(tableName, tableKey)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Query(selectQuery)
		suite.Nil(err)
	})
}

func (suite *FetcherDialectorSuite) TestSelectInfoQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()
	selectQuery := dialector.TableSelectInfoQuery(tableName)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Query(selectQuery)
		suite.Nil(err)
	})
}

func (suite *FetcherDialectorSuite) TestSelectSizeQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()
	selectQuery := dialector.TableSelectSizeQuery(tableName)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Query(selectQuery)
		suite.Nil(err)
	})
}

func (suite *FetcherDialectorSuite) TestSelectColumnsQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()
	selectQuery := dialector.TableSelectColumnsQuery(tableName)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Query(selectQuery)
		suite.Nil(err)
	})
}

func (suite *FetcherDialectorSuite) TestSelectIndexQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()
	selectQuery := dialector.TableSelectIndexesQuery(tableName)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Query(selectQuery)
		suite.Nil(err)
	})
}

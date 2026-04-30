package dialectors

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"journal-migrator/test"
	"testing"
)

func TestMigratorDialectorSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(MigratorDialectorSuite))
}

type MigratorDialectorSuite struct {
	test.IntegrationSuite
}

func (suite *MigratorDialectorSuite) createDialector() dialectors.BaseMigratorDialector {
	config := suite.Client.GetConfig()

	formatter := formatters.NewFormatter(config.Dialect)
	dialector := dialectors.NewMigratorDialector(formatter, config.Name)

	return dialector
}

func (suite *MigratorDialectorSuite) TestCreateColumn() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		var createQuery string
		var err error

		createQuery = dialector.CreateColumnQuery(tableName, "test_column_1", "VARCHAR(255)", "example", false)
		_, err = tx.Exec(createQuery)
		suite.Nil(err)
		createQuery = dialector.CreateColumnQuery(tableName, "test_column_2", "VARCHAR(255)", nil, false)
		_, err = tx.Exec(createQuery)
		suite.Nil(err)
		createQuery = dialector.CreateColumnQuery(tableName, "test_column_3", "BIGINT", 1, false)
		_, err = tx.Exec(createQuery)
		suite.Nil(err)
	})
}

func (suite *MigratorDialectorSuite) TestChangeColumn() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		var createQuery string
		var err error

		createQuery = dialector.ChangeColumnQuery(tableName, "data", "VARCHAR(255)", "example", false)
		_, err = tx.Exec(createQuery)
		suite.Nil(err)
		createQuery = dialector.ChangeColumnQuery(tableName, "data", "VARCHAR(255)", nil, false)
		_, err = tx.Exec(createQuery)
		suite.Nil(err)
		createQuery = dialector.ChangeColumnQuery(tableName, "data", "BIGINT", 1, false)
		_, err = tx.Exec(createQuery)
		suite.Nil(err)
	})
}

func (suite *MigratorDialectorSuite) TestRemoveColumn() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()
	removeQuery := dialector.RemoveColumnQuery(tableName, "data")

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Exec(removeQuery)
		suite.Nil(err)
	})
}

func (suite *MigratorDialectorSuite) TestRenameColumn() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()

	dialector := suite.createDialector()
	removeQuery := dialector.RenameColumnQuery(tableName, "data", "data_new")

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		_, err := tx.Exec(removeQuery)
		suite.Nil(err)
	})
}

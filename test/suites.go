package test

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/storage"
	"journal-migrator/test/helpers"
	"log/slog"
)

type IntegrationSuite struct {
	suite.Suite
	helpers.TableHelper

	Client *storage.DatabaseClient
	Logger *slog.Logger
}

func (suite *IntegrationSuite) SetupSuite() {
	suite.Client = SetupDatabaseClient("journal_migrator_test")
	suite.Logger = SetupLogger()
	suite.Config = suite.Client.GetConfig()

	err := suite.Client.OpenConnection()
	if err != nil {
		panic(err)
	}
}

func (suite *IntegrationSuite) TearDownSuite() {
	suite.Client.CloseConnection()
}

func (suite *IntegrationSuite) WithTable(tableName string, tableDef string, db storage.Database, testFn func(db storage.Database)) {
	tableCreateQuery := fmt.Sprintf("CREATE TABLE %s %s", tableName, tableDef)
	tableDeleteQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)

	_, _ = db.Exec(tableDeleteQuery)
	_, _ = db.Exec(tableCreateQuery)
	defer db.Exec(tableDeleteQuery)

	testFn(db)
}

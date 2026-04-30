package handlers

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/handlers"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/models/operation"
	"journal-migrator/pkg/models/operation/specs"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"journal-migrator/test"
	"testing"
)

func TestMigratorHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(MigratorHandlerSuite))
}

type MigratorHandlerSuite struct {
	test.IntegrationSuite
}

func (suite *MigratorHandlerSuite) createHandler(database storage.Database) *handlers.MigratorHandler {
	config := *suite.Client.GetConfig()

	formatter := formatters.NewFormatter(config.Dialect)
	dialector := dialectors.NewMigratorDialector(formatter, config.Name)
	handler := handlers.NewMigratorHandler(dialector, database, config, suite.Logger)

	return handler
}

func (suite *MigratorHandlerSuite) TestMigrateColumn() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()
	tableOps := []models.Operation{
		{
			Type: operation.CreateColumnType,
			Spec: &specs.CreateColumnSpec{
				Name:     "created_column",
				Type:     "INT",
				Default:  nil,
				Nullable: false,
			},
		},
		{
			Type: operation.ChangeColumnType,
			Spec: &specs.ChangeColumnSpec{
				Name:     "created_column",
				Type:     "BIGINT",
				Default:  nil,
				Nullable: false,
			},
		},
		{
			Type: operation.RenameColumnType,
			Spec: &specs.RenameColumnSpec{
				OldName: "created_column",
				NewName: "renamed_column",
			},
		},
		{
			Type: operation.RemoveColumnType,
			Spec: &specs.RemoveColumnSpec{
				Name: "renamed_column",
			},
		},
	}

	handler := suite.createHandler(tx)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		err := handler.Migrate(tableName, tableOps)
		suite.Nil(err)
	})
}

func (suite *MigratorHandlerSuite) TestMigrateIndex() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableName := "test_table"
	tableDef := suite.GetTableDefinition()
	tableOps := []models.Operation{
		{
			Type: operation.CreateIndexType,
			Spec: &specs.CreateIndexSpec{
				Name:    "created_index",
				Columns: []string{"id"},
				Unique:  true,
			},
		},
		{
			Type: operation.RenameIndexType,
			Spec: &specs.RenameIndexSpec{
				OldName: "created_index",
				NewName: "renamed_index",
			},
		},
		{
			Type: operation.RemoveIndexType,
			Spec: &specs.RemoveIndexSpec{
				Name: "renamed_index",
			},
		},
	}

	handler := suite.createHandler(tx)

	suite.WithTable(tableName, tableDef, tx, func(db storage.Database) {
		err := handler.Migrate(tableName, tableOps)
		suite.Nil(err)
	})
}

func (suite *MigratorHandlerSuite) TestMigrateTable() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	tableOps := []models.Operation{
		{
			Type: operation.CreateTableType,
			Spec: &specs.CreateTableSpec{
				Name:       "created_table",
				Definition: "id SERIAL PRIMARY KEY, data TEXT",
			},
		},
		{
			Type: operation.RenameTableType,
			Spec: &specs.RenameTableSpec{
				OldName: "created_table",
				NewName: "renamed_table",
			},
		},
		{
			Type: operation.RemoveTableType,
			Spec: &specs.RemoveTableSpec{
				Name: "renamed_table",
			},
		},
	}

	handler := suite.createHandler(tx)

	err := handler.Migrate("", tableOps)
	suite.Nil(err)
}

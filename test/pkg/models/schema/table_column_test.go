package schema

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/models"
	"testing"
)

func TestTableColumn(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping unit tests")
	}

	suite.Run(t, new(TableColumnSuite))
}

type TableColumnSuite struct {
	suite.Suite
}

func (suite *TableColumnSuite) TestIntersectColumns() {
	tableColumns_1 := []models.ColumnInfo{
		{Name: "column_1", Type: "VARCHAR(255)"},
		{Name: "column_2", Type: "VARCHAR(255)"},
	}

	tableColumns_2 := []models.ColumnInfo{
		{Name: "column_2", Type: "VARCHAR(255)"},
		{Name: "column_3", Type: "VARCHAR(255)"},
	}

	commonColumns := models.IntersectColumns(tableColumns_1, tableColumns_2)

	suite.Contains(commonColumns, "column_2")
	suite.NotContains(commonColumns, "column_1")
	suite.NotContains(commonColumns, "column_3")
}

package schema

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/models"
	"testing"
)

func TestTableIndex(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping unit tests")
	}

	suite.Run(t, new(TableIndexSuite))
}

type TableIndexSuite struct {
	suite.Suite
}

func (suite *TableIndexSuite) TestIntersectIndexes() {
	tableIndexes_1 := []models.IndexInfo{
		{Name: "index_1", Type: "BTREE"},
		{Name: "index_2", Type: "BTREE"},
	}

	tableIndexes_2 := []models.IndexInfo{
		{Name: "index_2", Type: "BTREE"},
		{Name: "index_3", Type: "BTREE"},
	}

	commonIndixes := models.IntersectIndexes(tableIndexes_1, tableIndexes_2)

	suite.Contains(commonIndixes, "index_2")
	suite.NotContains(commonIndixes, "index_1")
	suite.NotContains(commonIndixes, "index_3")
}

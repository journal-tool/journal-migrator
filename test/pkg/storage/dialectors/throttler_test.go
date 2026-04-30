package dialectors

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"journal-migrator/test"
	"testing"
)

func TestThrotterDialectorSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(ThrotterDialectorSuite))
}

type ThrotterDialectorSuite struct {
	test.IntegrationSuite
}

func (suite *ThrotterDialectorSuite) createDialector() dialectors.BaseThrottlerDialector {
	config := suite.Client.GetConfig()

	formatter := formatters.NewFormatter(config.Dialect)
	dialector := dialectors.NewThrottlerDialector(formatter)

	return dialector
}

func (suite *ThrotterDialectorSuite) TestReplicaHostsQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	dialector := suite.createDialector()
	selectQuery := dialector.SelectReplicaHostsQuery()

	_, err := tx.Exec(selectQuery)
	suite.Nil(err)
}

func (suite *ThrotterDialectorSuite) TestReplicaLagQuery() {
	db := suite.Client.GetDatabase()
	tx, _ := db.Begin()
	defer tx.Rollback()

	dialector := suite.createDialector()
	selectQuery := dialector.SelectReplicaLagQuery()

	_, err := tx.Exec(selectQuery)
	suite.Nil(err)
}

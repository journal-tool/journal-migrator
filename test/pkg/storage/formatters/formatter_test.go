package formatters

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/formatters"
	"strings"
	"testing"
)

func TestFormatterSuite(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping unit tests")
	}

	suite.Run(t, new(FormatterSuite))
}

type FormatterSuite struct {
	suite.Suite
}

func (suite *FormatterSuite) dialectors() []string {
	return []string{storage.MysqlDialect, storage.PostgresDialect}
}

func (suite *FormatterSuite) TestPrefixName() {
	for _, dialect := range suite.dialectors() {
		formatter := formatters.NewFormatter(dialect)
		maxLength := formatter.MaxLength()

		longName := strings.Repeat("a", maxLength)
		longName = formatter.PrefixName("prefix", longName)

		suite.True(strings.HasPrefix(longName, "prefix"))
		suite.True(maxLength <= len(longName))
	}
}

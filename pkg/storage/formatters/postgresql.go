package formatters

import (
	"fmt"
	"strings"
)

type PostgreSQLFormatter struct {
	maxLength int
}

func NewPostgreSQLFormatter() *PostgreSQLFormatter {
	return &PostgreSQLFormatter{
		maxLength: 63,
	}
}

func (f *PostgreSQLFormatter) MaxLength() int {
	return f.maxLength
}

func (f *PostgreSQLFormatter) PrefixName(prefix string, name string) string {
	prefixedName := fmt.Sprintf(`%s_%s`, prefix, name)
	truncatedLen := min(len(prefixedName), f.maxLength)
	return prefixedName[:truncatedLen]
}

func (f *PostgreSQLFormatter) ScopeEntity(scope string, name string) string {
	return fmt.Sprintf(`"%s"."%s"`, scope, name)
}

func (f *PostgreSQLFormatter) QuoteEntity(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

func (f *PostgreSQLFormatter) QuoteValue(value any) string {
	switch value.(type) {
	case string:
		return fmt.Sprintf(`'%v'`, value)
	default:
		return fmt.Sprintf(`%v`, value)
	}
}

func (f *PostgreSQLFormatter) TrimSpaces(query string) string {
	return strings.Join(strings.Fields(query), ` `)
}

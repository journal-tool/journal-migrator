package formatters

import (
	"fmt"
	"strings"
)

type MySQLFormatter struct {
	maxLength int
}

func NewMySQLFormatter() *MySQLFormatter {
	return &MySQLFormatter{
		maxLength: 64,
	}
}

func (f *MySQLFormatter) MaxLength() int {
	return f.maxLength
}

func (f *MySQLFormatter) PrefixName(prefix string, name string) string {
	prefixedName := fmt.Sprintf("%s_%s", prefix, name)
	truncatedLen := min(len(prefixedName), f.maxLength)
	return prefixedName[:truncatedLen]
}

func (f *MySQLFormatter) ScopeEntity(scope string, name string) string {
	return fmt.Sprintf("`%s`.`%s`", scope, name)
}

func (f *MySQLFormatter) QuoteEntity(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (f *MySQLFormatter) QuoteValue(value any) string {
	switch value.(type) {
	case string:
		return fmt.Sprintf("'%v'", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (f *MySQLFormatter) TrimSpaces(query string) string {
	return strings.Join(strings.Fields(query), " ")
}

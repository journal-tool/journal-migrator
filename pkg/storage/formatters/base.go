package formatters

import "journal-migrator/pkg/storage"

type BaseFormatter interface {
	MaxLength() int

	PrefixName(prefix string, name string) string
	ScopeEntity(scope string, name string) string
	QuoteEntity(name string) string
	QuoteValue(value any) string
	TrimSpaces(query string) string
}

func NewFormatter(dialect string) BaseFormatter {
	switch dialect {
	case storage.MysqlDialect:
		return NewMySQLFormatter()
	case storage.PostgresDialect:
		return NewPostgreSQLFormatter()
	}

	return nil
}

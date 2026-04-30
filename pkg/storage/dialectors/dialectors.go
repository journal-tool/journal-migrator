package dialectors

import (
	"journal-migrator/pkg/storage/dialectors/executor"
	"journal-migrator/pkg/storage/dialectors/fetcher"
	"journal-migrator/pkg/storage/dialectors/migrator"
	"journal-migrator/pkg/storage/dialectors/throttler"
	"journal-migrator/pkg/storage/formatters"
)

type BaseExecutorDialector = executor.BaseExecutorDialector
type BaseFetcherDialector = fetcher.BaseFetcherDialector
type BaseMigratorDialector = migrator.BaseMigratorDialector
type BaseThrottlerDialector = throttler.BaseThrottlerDialector

func NewExecutorDialector(formatter formatters.BaseFormatter, schema string) BaseExecutorDialector {
	switch format := formatter.(type) {
	case *formatters.MySQLFormatter:
		return executor.NewMySQLExecutorDialector(format, schema)
	case *formatters.PostgreSQLFormatter:
		return executor.NewPostgreSQLExecutorDialector(format, "public")
	}

	return nil
}

func NewFetcherDialector(formatter formatters.BaseFormatter, schema string) BaseFetcherDialector {
	switch format := formatter.(type) {
	case *formatters.MySQLFormatter:
		return fetcher.NewMySQLFetcherDialector(format, schema)
	case *formatters.PostgreSQLFormatter:
		return fetcher.NewPostgreSQLFetcherDialector(format, "public")
	}

	return nil
}

func NewMigratorDialector(formatter formatters.BaseFormatter, schema string) BaseMigratorDialector {
	switch format := formatter.(type) {
	case *formatters.MySQLFormatter:
		return migrator.NewMySQLMigratorDialector(format, schema)
	case *formatters.PostgreSQLFormatter:
		return migrator.NewPostgreSQLMigratorDialector(format, "public")
	}

	return nil
}

func NewThrottlerDialector(formatter formatters.BaseFormatter) BaseThrottlerDialector {
	switch format := formatter.(type) {
	case *formatters.MySQLFormatter:
		return throttler.NewMySQLThrottlerDialector(format)
	case *formatters.PostgreSQLFormatter:
		return throttler.NewPostgreSQLThrottlerDialector(format)
	}

	return nil
}

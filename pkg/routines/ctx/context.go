package ctx

import (
	"journal-migrator/pkg/handlers"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"journal-migrator/pkg/storage/formatters"
	"log/slog"
)

type Context struct {
	config *storage.DatabaseConfig
	logger *slog.Logger
}

func NewContext(config *storage.DatabaseConfig, logger *slog.Logger) *Context {
	return &Context{
		config: config,
		logger: logger,
	}
}

func (ctx *Context) CreateExecutor(database storage.Database) *handlers.ExecutorHandler {
	formatter := formatters.NewFormatter(ctx.config.Dialect)
	dialector := dialectors.NewExecutorDialector(formatter, ctx.config.Name)
	handler := handlers.NewExecutorHandler(dialector, database, *ctx.config, ctx.logger)

	return handler
}

func (ctx *Context) CreateFetcher(database storage.Database) *handlers.FetcherHandler {
	formatter := formatters.NewFormatter(ctx.config.Dialect)
	dialector := dialectors.NewFetcherDialector(formatter, ctx.config.Name)
	handler := handlers.NewFetcherHandler(dialector, database, *ctx.config, ctx.logger)

	return handler
}

func (ctx *Context) CreateMigrator(database storage.Database) *handlers.MigratorHandler {
	formatter := formatters.NewFormatter(ctx.config.Dialect)
	dialector := dialectors.NewMigratorDialector(formatter, ctx.config.Name)
	handler := handlers.NewMigratorHandler(dialector, database, *ctx.config, ctx.logger)

	return handler
}

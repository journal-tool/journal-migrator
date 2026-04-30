package throttler

import "journal-migrator/pkg/storage/formatters"

type PostgreSQLThrottlerDialector struct {
	formatter *formatters.PostgreSQLFormatter
}

func NewPostgreSQLThrottlerDialector(formatter *formatters.PostgreSQLFormatter) *PostgreSQLThrottlerDialector {
	return &PostgreSQLThrottlerDialector{
		formatter: formatter,
	}
}

func (d *PostgreSQLThrottlerDialector) SelectReplicaHostsQuery() string {
	return d.formatter.TrimSpaces(`
		SELECT
			client_addr,
			client_port
		FROM
			pg_catalog.pg_stat_replication
	`)
}

func (d *PostgreSQLThrottlerDialector) SelectReplicaLagQuery() string {
	return d.formatter.TrimSpaces(`
		SELECT
			EXTRACT(
				SECOND FROM (
					NOW() - COALESCE(pg_last_xact_replay_timestamp(), NOW())
				)
			)
	`)
}

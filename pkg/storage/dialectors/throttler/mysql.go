package throttler

import "journal-migrator/pkg/storage/formatters"

type MySQLThrottlerDialector struct {
	formatter *formatters.MySQLFormatter
}

func NewMySQLThrottlerDialector(formatter *formatters.MySQLFormatter) *MySQLThrottlerDialector {
	return &MySQLThrottlerDialector{
		formatter: formatter,
	}
}

func (d *MySQLThrottlerDialector) SelectReplicaHostsQuery() string {
	return d.formatter.TrimSpaces(`
		SELECT
			SUBSTRING_INDEX(host, ':', +1),
			SUBSTRING_INDEX(host, ':', -1)
		FROM
			information_schema.processlist
		WHERE
			command = 'Binlog Dump'
	`)
}

func (d *MySQLThrottlerDialector) SelectReplicaLagQuery() string {
	return d.formatter.TrimSpaces(`
		SELECT
			processlist_time
		FROM
			performance_schema.threads
		JOIN
			performance_schema.replication_applier_status_by_worker USING (thread_id)
		WHERE
			channel_name = 'group_replication_applier'
	`)
}

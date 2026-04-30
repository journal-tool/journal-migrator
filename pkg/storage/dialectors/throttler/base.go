package throttler

type BaseThrottlerDialector interface {
	SelectReplicaHostsQuery() string
	SelectReplicaLagQuery() string
}

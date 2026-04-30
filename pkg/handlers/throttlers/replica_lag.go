package throttlers

import (
	"journal-migrator/lib"
	"journal-migrator/pkg/models"
	"journal-migrator/pkg/storage"
	"journal-migrator/pkg/storage/dialectors"
	"log/slog"
	"slices"
	"sync"
	"time"
)

const defaultWaitingSecs = 00
const maximumWaitingSecs = 60
const replicaLagCheckPeriod = time.Minute

var ParseReplicationHostRows = lib.ParseRows[models.ReplicationHost]
var ParseReplicationLagRow = lib.ParseRow[models.ReplicationLag]

type ReplicaLagThrottler struct {
	dialector     dialectors.BaseThrottlerDialector
	database      storage.Database
	config        storage.DatabaseConfig
	logger        *slog.Logger
	lastTimeCheck time.Time
	nextTimeCheck time.Time
	waitingSecs   int
}

func NewReplicaLagThrottler(dialector dialectors.BaseThrottlerDialector, database storage.Database, config storage.DatabaseConfig, logger *slog.Logger) *ReplicaLagThrottler {
	return &ReplicaLagThrottler{
		dialector:     dialector,
		database:      database,
		config:        config,
		logger:        logger,
		lastTimeCheck: time.Time{},
		nextTimeCheck: time.Time{},
		waitingSecs:   defaultWaitingSecs,
	}
}

func (t *ReplicaLagThrottler) buildReplicaClient(host models.ReplicationHost) (*storage.DatabaseClient, error) {
	config := t.config
	config.Host = host.Host
	config.Port = host.Port

	client := storage.NewDatabaseClient(&config)

	err := client.OpenConnection()
	if err != nil {
		t.logger.Error("Cannot connect to replica",
			slog.String("reason", err.Error()),
			slog.String("db_host", config.Host),
			slog.String("db_port", config.Port),
		)
		return nil, err
	}

	return client, nil
}

func (t *ReplicaLagThrottler) fetchReplicaHosts() ([]models.ReplicationHost, error) {
	query := t.dialector.SelectReplicaHostsQuery()

	rows, err := t.database.Query(query)
	if err != nil {
		t.logger.Error("Cannot query source replication info",
			slog.String("reason", err.Error()),
			slog.String("db_host", t.config.Host),
			slog.String("db_port", t.config.Port),
		)
		return nil, err
	}

	hosts, err := ParseReplicationHostRows(rows)
	if err != nil {
		t.logger.Error("Cannot parse source replication info",
			slog.String("reason", err.Error()),
			slog.String("db_host", t.config.Host),
			slog.String("db_port", t.config.Port),
		)
	}

	return hosts, nil
}

func (t *ReplicaLagThrottler) fetchReplicaLag(host models.ReplicationHost) (int, error) {
	client, err := t.buildReplicaClient(host)
	if err != nil {
		return t.waitingSecs, err
	}

	defer client.CloseConnection()

	query := t.dialector.SelectReplicaLagQuery()
	database := client.GetDatabase()

	row := database.QueryRow(query)
	err = row.Err()
	if err != nil {
		t.logger.Warn("Cannot query replica lag",
			slog.String("reason", err.Error()),
			slog.String("db_host", t.config.Host),
			slog.String("db_port", t.config.Port),
		)
		return t.waitingSecs, err
	}

	replicaLag, err := ParseReplicationLagRow(row)
	if err != nil {
		t.logger.Warn("Cannot parse replica lag",
			slog.String("reason", err.Error()),
			slog.String("db_host", t.config.Host),
			slog.String("db_port", t.config.Port),
		)
		return t.waitingSecs, err
	}

	return replicaLag.Seconds, nil
}

func (t *ReplicaLagThrottler) Throttle() {
	timeNow := time.Now().UTC()
	if timeNow.Before(t.nextTimeCheck) {
		return
	}

	hosts, err := t.fetchReplicaHosts()
	if err != nil {
		return
	}

	var lagSeconds = make([]int, len(hosts))
	var wg sync.WaitGroup

	for i, host := range hosts {
		if host.IsSource() {
			continue
		}
		wg.Go(func() {
			lagSeconds[i], _ = t.fetchReplicaLag(host)
		})
	}

	wg.Wait()

	maxReplicaLagSecs := slices.Max(lagSeconds)
	t.lastTimeCheck = timeNow
	t.nextTimeCheck = timeNow.Add(replicaLagCheckPeriod)
	t.waitingSecs = min(maxReplicaLagSecs, maximumWaitingSecs)

	time.Sleep(
		time.Duration(t.waitingSecs) * time.Second,
	)
}

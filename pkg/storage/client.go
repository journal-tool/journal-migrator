package storage

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"time"
)

type DatabaseClient struct {
	config *DatabaseConfig
	conn   *sql.DB
}

func NewDatabaseClient(config *DatabaseConfig) *DatabaseClient {
	return &DatabaseClient{
		config: config,
	}
}

func (client *DatabaseClient) openConnection(dialect string, connectionDSN string) error {
	var currentBackoff = 02 * time.Second
	var maximumBackoff = 64 * time.Second

	for currentBackoff < maximumBackoff {
		conn, err := sql.Open(dialect, connectionDSN)
		if err != nil {
			return err
		}

		err = conn.Ping()
		if err != nil {
			time.Sleep(currentBackoff)
			currentBackoff *= 2
			continue
		}

		client.conn = conn
		return nil
	}

	return errors.New("cannot connect to the database")
}

func (client *DatabaseClient) GetConfig() *DatabaseConfig {
	return client.config
}

func (client *DatabaseClient) GetDatabase() *sql.DB {
	return client.conn
}

func (client *DatabaseClient) OpenConnection() error {
	dsn, err := client.config.BuildConnectionDSN()
	if err != nil {
		return err
	}

	return client.openConnection(client.config.Dialect, dsn)
}

func (client *DatabaseClient) CloseConnection() {
	if client.conn != nil {
		_ = client.conn.Close()
	}
}

package storage

import (
	"errors"
	"fmt"
)

const MysqlDialect = "mysql"
const PostgresDialect = "postgres"

type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	Dialect  string
}

func NewDatabaseConfig(user string, pass string, host string, port string, name string, dialect string) *DatabaseConfig {
	return &DatabaseConfig{
		Username: user,
		Password: pass,
		Host:     host,
		Port:     port,
		Name:     name,
		Dialect:  dialect,
	}
}

func (config *DatabaseConfig) buildMysqlDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}

func (config *DatabaseConfig) buildPostgresDSN() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}

func (config *DatabaseConfig) BuildConnectionDSN() (string, error) {
	var connectionDSN string

	switch config.Dialect {
	case MysqlDialect:
		connectionDSN = config.buildMysqlDSN()
	case PostgresDialect:
		connectionDSN = config.buildPostgresDSN()
	default:
		return "", errors.New("invalid database dialect")
	}

	return connectionDSN, nil
}

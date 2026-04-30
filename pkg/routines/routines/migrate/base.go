package migrate

import (
	"journal-migrator/pkg/handlers/throttlers"
	"journal-migrator/pkg/models"
)

type BaseMigrateRoutine interface {
	Progress() int64
	Run(table string, operations []models.Operation, throttler throttlers.BaseThrottler) error
}

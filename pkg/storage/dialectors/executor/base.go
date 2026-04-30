package executor

import "journal-migrator/pkg/models"

type BaseExecutorDialector interface {
	TableCreateQuery(sourceTable string, targetTable string) string
	TableRemoveQuery(sourceTable string) string

	TableInsertBatchQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string, lowerID int, upperID int) string

	TriggerCreateDeleteQuery(sourceTable models.TableInfo, targetTable models.TableInfo) string
	TriggerCreateInsertQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string) string
	TriggerCreateUpdateQuery(sourceTable models.TableInfo, targetTable models.TableInfo, columns []string) string

	TriggerDropDeleteQuery(sourceTable string) string
	TriggerDropInsertQuery(sourceTable string) string
	TriggerDropUpdateQuery(sourceTable string) string
}

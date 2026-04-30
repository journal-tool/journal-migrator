package models

import (
	"journal-migrator/pkg/models/operation"
	"journal-migrator/pkg/models/schema"
	"journal-migrator/pkg/models/throttler"
)

type Operation = operation.Operation
type ColumnInfo = schema.ColumnInfo
type IndexInfo = schema.IndexInfo
type TableInfo = schema.TableInfo
type TableSize = schema.TableSize
type ReplicationHost = throttler.ReplicationHost
type ReplicationLag = throttler.ReplicationLag

var ParseOperations = operation.ParseOperations
var IntersectColumns = schema.IntersectColumns
var IntersectIndexes = schema.IntersectIndexes

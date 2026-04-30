package cleanup

type BaseCleanupRoutine interface {
	Progress() int64
	Run(table string) error
}

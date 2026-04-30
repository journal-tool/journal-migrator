package throttlers

type BaseThrottler interface {
	Throttle()
}

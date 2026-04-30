package throttlers

import "time"

type WaitingTimeThrottler struct {
	waitingSecs int
}

func NewWaitingTimeThrottler(waitSeconds int) *WaitingTimeThrottler {
	return &WaitingTimeThrottler{
		waitingSecs: waitSeconds,
	}
}

func (t *WaitingTimeThrottler) Throttle() {
	time.Sleep(
		time.Duration(t.waitingSecs) * time.Second,
	)
}

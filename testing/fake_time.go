package testing

import (
	"time"

	"github.com/go-co-op/gocron"
)

type FakeTime struct {
	NowTime time.Time
}

var _ gocron.TimeWrapper = &FakeTime{}

func NewFakeTime(now time.Time) *FakeTime {
	return &FakeTime{
		NowTime: now,
	}
}

func (ft *FakeTime) Now(_ *time.Location) time.Time {
	return ft.NowTime
}

func (ft *FakeTime) Unix(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func (ft *FakeTime) Sleep(d time.Duration) {
	if d < 0 {
		return
	}
	ft.NowTime = ft.NowTime.Add(d)
}

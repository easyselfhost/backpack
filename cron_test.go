package backpack_test

import (
	"testing"
	"time"

	bp "github.com/easyselfhost/backpack"
	bt "github.com/easyselfhost/backpack/testing"
	"github.com/go-co-op/gocron"
	"github.com/golang/mock/gomock"
)

var (
	testLocation = time.UTC
	testTime     = time.Date(2023, 1, 1, 0, 0, 0, 0, testLocation)
)

type cronFields struct {
	elapsed []time.Duration
}

type cronArgs struct {
	name  string
	daily []string
	every string
}

type cronTestCase struct {
	name    string
	fields  cronFields
	args    cronArgs
	wantErr bool
}

func cronTest(tt *cronTestCase) func(*testing.T) {
	return func(t *testing.T) {
		ft := bt.NewFakeTime(testTime)
		scheduler := gocron.NewScheduler(testLocation)
		scheduler.CustomTime(ft)
		scheduler.CustomTimer(func(d time.Duration, f func()) *time.Timer {
			if d > 0 {
				ft.NowTime = ft.NowTime.Add(d)
			}
			return time.AfterFunc(d*0, f)
		})
		c := bp.NewCronWithScheduler(scheduler)

		ctrl := gomock.NewController(t)
		wf := bt.NewMockWorkflow(ctrl)

		var err error
		if len(tt.args.daily) > 0 {
			err = c.RegisterDaily(tt.args.name, wf, tt.args.daily)
		} else {
			err = c.RegisterInterval(tt.args.name, wf, tt.args.every)
		}
		if (err != nil) != tt.wantErr {
			t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
		}
		if err != nil {
			return
		}

		var (
			jobStart = make(chan bool)
			jobEnd   = make(chan bool)
		)

		wf.EXPECT().Run().DoAndReturn(func() error {
			jobStart <- true
			<-jobEnd
			return nil
		}).AnyTimes()

		c.StartAsync()

		curr := testTime
		for _, d := range tt.fields.elapsed {
			<-jobStart
			curr = curr.Add(d)
			if curr.After(ft.NowTime) {
				t.Error("expected time after execution is after now time")
			}
			jobEnd <- true
		}
	}
}

func TestCron_RegisterDaily(t *testing.T) {
	tests := []cronTestCase{
		{
			name: "invalid time format",
			args: cronArgs{
				name:  "daily-job",
				daily: []string{"a", "b"},
			},
			wantErr: true,
		},
		{
			name: "single daily time single day",
			fields: cronFields{
				elapsed: []time.Duration{time.Hour * 3},
			},
			args: cronArgs{
				name:  "daily-job",
				daily: []string{"3:00"},
			},
			wantErr: false,
		},
		{
			name: "single daily time multiple days",
			fields: cronFields{
				elapsed: []time.Duration{
					time.Hour * 3,
					time.Hour * 24,
				},
			},
			args: cronArgs{
				name:  "daily-job",
				daily: []string{"3:00"},
			},
			wantErr: false,
		},
		{
			name: "multiple daily time",
			fields: cronFields{
				elapsed: []time.Duration{
					time.Hour * 3,
					time.Hour * 9,
					time.Hour * 7,
				},
			},
			args: cronArgs{
				name:  "daily-job",
				daily: []string{"3:00", "12:00", "19:00"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, cronTest(&tt))
	}
}

func TestCron_RegisterHourly(t *testing.T) {
	tests := []cronTestCase{
		{
			name: "invalid time format 1",
			args: cronArgs{
				name:  "hourly-job",
				every: "1t",
			},
			wantErr: true,
		},
		{
			name: "invalid time format 2",
			args: cronArgs{
				name:  "hourly-job",
				every: "3",
			},
			wantErr: true,
		},
		{
			name: "every hour",
			fields: cronFields{
				elapsed: []time.Duration{time.Hour, time.Hour},
			},
			args: cronArgs{
				name:  "hourly-job",
				every: "1h",
			},
			wantErr: false,
		},
		{
			name: "every 30 minutes",
			fields: cronFields{
				elapsed: []time.Duration{
					time.Minute * 30,
					time.Minute * 30,
				},
			},
			args: cronArgs{
				name:  "hourly-job",
				every: "30m",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, cronTest(&tt))
	}
}

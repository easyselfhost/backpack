package backpack

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/golang/glog"
)

type Cron struct {
	scheduler *gocron.Scheduler
	jobs      map[string]*gocron.Job
}

func NewCron() *Cron {
	return NewCronWithScheduler(gocron.NewScheduler(time.Local))
}

func NewCronWithScheduler(scheduler *gocron.Scheduler) *Cron {
	return &Cron{
		scheduler: scheduler,
		jobs:      make(map[string]*gocron.Job),
	}
}

func NewCronFromConfig(config *Config) (*Cron, error) {
	cron := NewCron()

	for i, rule := range config.BackupRules {
		schedule := &rule.Schedule
		if len(schedule.DailySchedule) == 0 && len(schedule.EveryInterval) == 0 {
			return nil, fmt.Errorf("no daily schedule or interval schedule is set")
		}

		wf := NewBackpackFlow(rule)

		if len(schedule.DailySchedule) > 0 {
			n := fmt.Sprintf("backup-%d-daily", i)
			err := cron.RegisterDaily(n, wf, schedule.DailySchedule)
			if err != nil {
				return nil, err
			}
		}

		if len(schedule.EveryInterval) > 0 {
			n := fmt.Sprintf("backup-%d-%s", i, schedule.EveryInterval)
			err := cron.RegisterInterval(n, wf, schedule.EveryInterval)
			if err != nil {
				return nil, err
			}
		}
	}

	return cron, nil
}

func genJobFunc(name string, wf Workflow) func() {
	return func() {
		err := wf.Run()
		if err != nil {
			glog.Errorf("Job %s: failed to run backup workflow %v", name, err)
		}
	}
}

func (c *Cron) StartAsync() {
	c.scheduler.StartAsync()
}

func (c *Cron) Stop() {
	c.scheduler.Stop()
	c.scheduler.Clear()
	c.jobs = make(map[string]*gocron.Job)
}

func (c *Cron) RegisterDaily(name string, wf Workflow, daily []string) error {
	if _, found := c.jobs[name]; found {
		return errors.New("job already exists")
	}

	dailyTimes := strings.Join(daily, ";")
	dailyJob, err := c.scheduler.Every(1).Day().At(dailyTimes).Do(genJobFunc(name, wf))
	if err != nil {
		return fmt.Errorf("failed to create daily job: %w", err)
	}

	c.jobs[name] = dailyJob
	return nil
}

func (c *Cron) RegisterInterval(name string, wf Workflow, every string) error {
	if _, found := c.jobs[name]; found {
		return errors.New("job already exists")
	}

	hourlyJob, err := c.scheduler.Every(every).Do(genJobFunc(name, wf))
	if err != nil {
		return fmt.Errorf("failed to create interval job: %w", err)
	}

	c.jobs[name] = hourlyJob

	return nil
}

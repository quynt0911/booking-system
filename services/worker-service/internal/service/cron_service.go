package scheduler

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/your-org/booking-system/services/worker-service/internal/jobs"
)

type Cron struct {
	scheduler *cron.Cron
	jobs      map[string]jobs.Job
}

func NewCron() *Cron {
	return &Cron{
		scheduler: cron.New(),
		jobs:      make(map[string]jobs.Job),
	}
}

func (c *Cron) RegisterJob(name string, job jobs.Job) {
	c.jobs[name] = job
}

func (c *Cron) ScheduleJob(name, schedule string, payload interface{}) error {
	job, exists := c.jobs[name]
	if !exists {
		return fmt.Errorf("job %s not found", name)
	}

	_, err := c.scheduler.AddFunc(schedule, func() {
		if err := job.Execute(context.Background(), payload); err != nil {
			log.Printf("Error executing job %s: %v", name, err)
		}
	})

	return err
}

func (c *Cron) Start() {
	c.scheduler.Start()
}

func (c *Cron) Stop() {
	c.scheduler.Stop()
}

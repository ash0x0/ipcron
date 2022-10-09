package ipcron

import (
	"testing"
	"time"
)

func TestAddJobToSchedule(t *testing.T) {
	interval, _ := time.ParseDuration("1s")

	schedule := createTestSchedule()

	firstJob, _ := schedule.ScheduleJobWithInterval(interval, simpleJob, "firstJob")
	secondJob, _ := schedule.ScheduleJobWithInterval(interval, simpleJob, "secondJob")

	schedule.addJobToSchedule(firstJob)
	schedule.addJobToSchedule(secondJob)

	got := len(schedule.jobQueue)
	if got < 2 {
		t.Errorf("Not all jobs are added to schedule. Expected %v, Got %v.", 2, got)
	}

	name := schedule.jobQueue[0].name
	if name != "firstJob" {
		t.Errorf("First job not added to schedule. Expected %v, Got %v.", "firstJob", name)
	}
}

func TestSortSchedule(t *testing.T) {
	earlyInterval, _ := time.ParseDuration("1s")
	lateInterval, _ := time.ParseDuration("10s")

	schedule := createTestSchedule()

	earlyJob, _ := schedule.ScheduleJobWithInterval(earlyInterval, simpleJob, "simpleJob-early")
	lateJob, _ := schedule.ScheduleJobWithInterval(lateInterval, simpleJob, "simpleJob-late")

	schedule.addJobToSchedule(earlyJob)
	schedule.addJobToSchedule(lateJob)

	schedule.sortSchedule()

	if schedule.jobQueue[0].nextTime.Before(schedule.jobQueue[1].nextTime) {
		t.Errorf("Schedule isn't ordered correctly. Earliest job not at the top.")
	}
}

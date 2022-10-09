package ipcron

import (
	"testing"
)

func TestAddJobToSchedule(t *testing.T) {
	schedule := createTestSchedule()

	firstJob, _ := schedule.ScheduleJobWithInterval("1s", "1s", simpleJob, "firstJob")
	secondJob, _ := schedule.ScheduleJobWithInterval("1s", "1s", simpleJob, "secondJob")

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
	schedule := createTestSchedule()

	earlyJob, _ := schedule.ScheduleJobWithInterval("1s", "1s", simpleJob, name+"-early")
	lateJob, _ := schedule.ScheduleJobWithInterval("10s", "1s", simpleJob, name+"-late")

	schedule.addJobToSchedule(earlyJob)
	schedule.addJobToSchedule(lateJob)

	schedule.sortSchedule()

	if schedule.jobQueue[0].nextTime.Before(schedule.jobQueue[1].nextTime) {
		t.Errorf("Schedule isn't ordered correctly. Earliest job not at the top.")
	}
}

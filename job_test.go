package ipcron

import (
	"testing"
	"time"
)

func TestCreateSimpleJob(t *testing.T) {
	job, err := createTestSchedule().ScheduleJobWithInterval(interval, simpleJob, name)

	if err != nil {
		t.Errorf("Job creation error %v", err)
	}

	if job.name != name {
		t.Errorf("Job has wrong name. Expected %v, Got %v", name, job.name)
	}
	if len(job.id) < 1 {
		t.Errorf("Job is created without an ID")
	}
	if job.interval != interval {
		t.Errorf("Job has incorrect interval. Expected %v, Got %v", interval, job.interval)
	}
	if job.execCount != 0 {
		t.Errorf("Job starts with wrong exec count. Expected %v, Got %v", 0, job.execCount)
	}
}

func TestCreateCronExprJob(t *testing.T) {
	cronExpression := "* * * * * * *"
	job, err := createTestSchedule().ScheduleWithCronSyntax(cronExpression, simpleJob, name)

	if err != nil {
		t.Errorf("Cron expression job creation error %v", err)
	}

	if job.interval != interval {
		t.Errorf("Cron expression job has incorrect interval. Expected %v, Got %v", interval, job.interval)
	}
}

func TestSetExecLimit(t *testing.T) {
	execLimit := 5

	job, _ := createTestSchedule().ScheduleJobWithInterval(interval, simpleJob, name)
	job.SetExecutionLimit(execLimit)

	if job.execLimit != execLimit {
		t.Errorf("Cron expression job has incorrect interval. Expected %v, Got %v", execLimit, job.execLimit)
	}
}

func TestUpdateJob(t *testing.T) {
	job, _ := createTestSchedule().ScheduleJobWithInterval(interval, simpleJob, name)

	time.Sleep(2 * time.Second)

	job.updateJob()

	if !job.nextTime.After(time.Now()) {
		t.Errorf("Job update produces time in the past. Got %v", job.nextTime)
	}
}

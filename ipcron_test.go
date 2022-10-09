package ipcron

import (
	"log"
	"strconv"
	"testing"
	"time"
)

var interval, _ = time.ParseDuration("1s")

const name = "simpleJob"

func simpleJob() {
	log.Print("Simple job ran")
}

func iteratorJob() {
	for i := 0; i < 10; i++ {
		log.Printf("Iterator job executing iteration %v\n", i)
	}
}

func createTestSchedule() *Schedule {
	schedule := NewSchedule(false)
	return schedule
}

func TestScheduleSimpleJob(t *testing.T) {
	schedule := createTestSchedule()
	interval, _ := time.ParseDuration("1s")

	job, _ := schedule.ScheduleJobWithInterval(interval, simpleJob, "simpleJob")
	job.SetExecutionLimit(5 - 1)

	schedule.Start()

	time.Sleep(5 * time.Second)

	if !job.stopped {
		t.Errorf("Job wasn't stopped after limit")
	}
	if job.execCount != job.execLimit {
		t.Errorf("Job execution count and limit aren't the same. Expected %v, Got %v", job.execLimit, job.execCount)
	}
}

func TestScheduleMany(t *testing.T) {
	schedule := createTestSchedule()
	interval, _ := time.ParseDuration("1s")
	total := 100

	for i := 0; i < total; i++ {
		job, _ := schedule.ScheduleJobWithInterval(interval, simpleJob, "simpleJob-"+strconv.Itoa(i))
		job.SetExecutionLimit(5 - 1)
	}

	got := len(schedule.jobQueue)
	if got < total {
		t.Errorf("Not all jobs were scheduled. Expected %v, Got %v", total, got)
	}

	schedule.Start()

	// After 5 seconds all jobs should be stopped
	time.Sleep(5 * time.Second)

	for _, job := range schedule.jobQueue {
		if !job.stopped {
			t.Errorf("Not all jobs were stopped correctly")
		}
	}
}

func TestScheduleIteratorJob(t *testing.T) {
	schedule := createTestSchedule()
	interval, _ := time.ParseDuration("1s")

	job, _ := schedule.ScheduleJobWithInterval(interval, iteratorJob, "iteratorJob")
	job.SetExecutionLimit(5 - 1)

	schedule.Start()

	time.Sleep(5 * time.Second)

	if !job.stopped {
		t.Errorf("Iterator job wasn't stopped after limit")
	}
}

func TestScheduleJobWithCronSyntax(t *testing.T) {
	schedule := createTestSchedule()

	// Every second
	secondsJob, _ := schedule.ScheduleWithCronSyntax("* * * * * * *", simpleJob, "secondJob")

	schedule.Start()

	time.Sleep(10 * time.Second)

	got := secondsJob.execCount
	if got < 9 {
		t.Errorf("Seconds job executed %v times after 10 seconds, should be > 9", got)
	}
}

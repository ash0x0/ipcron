package main

import (
	"fmt"
	"ipcron"
	"time"
)

func tenSeconds() {
	fmt.Println("Every 10 seconds")
}

func everySecond() {
	fmt.Println("Every 1 second")
}

func main() {
	schedule := ipcron.NewSchedule(true)

	// Can create a new job by setting a time interval between each occurance
	// Receives intervals as time.Duration
	interval, _ := time.ParseDuration("10s")
	intervalJob, _ := schedule.ScheduleJobWithInterval(interval, tenSeconds, "intervalExample")

	// Can also use Cron Expression syntax to schedule the job at a certain time occurance
	expressionJob, _ := schedule.ScheduleWithCronSyntax("* * * * * * *", everySecond, "cronExpressionExample")

	// Can set the maximumum number of occurances for the job, for both cron syntax and interval jobs
	intervalJob.SetExecutionLimit(2)
	expressionJob.SetExecutionLimit(5)

	// After all jobs are added, need to start the scheduler
	schedule.Start()

	// Wait for a bit to see the jobs running
	time.Sleep((10 + 1) * time.Second)

	// Stop the schedule to kill all jobs and release any wait or lock
	schedule.Stop()
}

package ipcron

import (
	"errors"
	"log"
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

func (s *Schedule) ScheduleJobWithInterval(timeInterval string, maxRunTime string, job func(), jobName string) (*Job, error) {
	log.Printf("Adding a new job with interval to schedule. Name %v, Interval %v, run time %v\n", jobName, timeInterval, maxRunTime)

	parsedInterval, _ := time.ParseDuration(timeInterval)
	nextTime := time.Now().Add(parsedInterval)
	parsedRunTime, _ := time.ParseDuration(maxRunTime)

	newJob := Job{name: jobName, job: job, scheduleExpression: "", nextTime: nextTime, interval: parsedInterval, maxRunTime: parsedRunTime}

	s.addJobToSchedule(&newJob)

	return &newJob, nil
}

func (s *Schedule) ScheduleWithCronSyntax(scheduleExpression string, maxRunTime string, job func(), jobName string) (*Job, error) {
	log.Printf("Adding a new job with cron expr to schedule. Name %v, expression %v, run time %v\n", jobName, scheduleExpression, maxRunTime)
	var nextTime time.Time
	var timeInterval time.Duration

	nextTimeSlice := cronexpr.MustParse(scheduleExpression).NextN(time.Now(), 2)

	if len(nextTimeSlice) == 2 {
		nextTime = nextTimeSlice[0]
		timeInterval = nextTimeSlice[1].Sub(nextTime)
	} else {
		return new(Job), errors.New("provided schedule expression is invalid")
	}

	parsedRunTime, _ := time.ParseDuration(maxRunTime)

	newJob := Job{name: jobName, job: job, scheduleExpression: scheduleExpression, nextTime: nextTime, interval: timeInterval, maxRunTime: parsedRunTime}

	s.addJobToSchedule(&newJob)

	return &newJob, nil
}

func (s *Schedule) Stop() {
	log.Printf("Stopping scheduler with number of jobs %v at time %v\n", len(s.jobQueue), time.Now())
	for i := 0; i < len(s.jobQueue); i++ {
		if !s.jobQueue[i].stopped {
			s.jobWaitGroup.Done()
			log.Printf("Stopped job %v with ID %v at time %v\n", s.jobQueue[i].name, s.jobQueue[i].id, time.Now())
		}
	}
	log.Printf("Stopped all jobs at time %v\n", time.Now())
	s.scheduleWaitGroup.Done()
	log.Printf("Stopped main scheduler rouitine and closed log at time %v\n", time.Now())
	s.logFile.Close()
}

func (s *Schedule) Start() {
	s.scheduleWaitGroup.Add(1)
	log.Printf("Starting scheduler with number of jobs %v at time %v\n", len(s.jobQueue), time.Now())
	go s.startJobRoutines()
}

func NewSchedule(enableLogging bool) *Schedule {
	log.SetPrefix("ipcron: ")
	log.SetFlags(0)

	var jobQueue []*Job
	schedule := &Schedule{jobQueue: jobQueue}
	if enableLogging {
		schedule.openLog()
	}
	log.Printf("Created new scheduler at time %v\n", time.Now())
	return schedule
}

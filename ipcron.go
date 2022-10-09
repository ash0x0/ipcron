package ipcron

import (
	"errors"
	"log"
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

func (s *Schedule) ScheduleJobWithInterval(timeInterval time.Duration, job func(), jobName string) (*Job, error) {
	nextTime := time.Now().Add(timeInterval)
	newJob := Job{name: jobName, job: job, scheduleExpression: "", nextTime: nextTime, interval: timeInterval}
	s.addJobToSchedule(&newJob)
	return &newJob, nil
}

func (s *Schedule) ScheduleWithCronSyntax(scheduleExpression string, job func(), jobName string) (*Job, error) {
	var nextTime time.Time
	var timeInterval time.Duration
	nextTimeSlice := cronexpr.MustParse(scheduleExpression).NextN(time.Now(), 2)
	if len(nextTimeSlice) == 2 {
		nextTime = nextTimeSlice[0]
		timeInterval = nextTimeSlice[1].Sub(nextTime)
	} else {
		return new(Job), errors.New("provided schedule expression is invalid")
	}
	newJob := Job{name: jobName, job: job, scheduleExpression: scheduleExpression, nextTime: nextTime, interval: timeInterval}
	s.addJobToSchedule(&newJob)
	return &newJob, nil
}

func (s *Schedule) Stop() {
	for i := 0; i < len(s.jobQueue); i++ {
		if !s.jobQueue[i].stopped {
			s.jobWaitGroup.Done()
		}
	}
	s.scheduleWaitGroup.Done()
	s.logFile.Close()
}

func (s *Schedule) Start() {
	s.scheduleWaitGroup.Add(1)
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
	return schedule
}

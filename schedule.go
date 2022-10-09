package ipcron

import (
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type Schedule struct {
	jobQueue          []*Job
	jobWaitGroup      sync.WaitGroup
	scheduleWaitGroup sync.WaitGroup
	logFile           *os.File
}

func (s *Schedule) startJob(j *Job) {
	log.Printf("Starting job %v routine at time %v\n", j.name, time.Now())
	j.updateJob()
	for j.nextTime.After(time.Now()) {
		// This is because I'm yet unsure of Go's short circuit behavior
		if j.execLimit != 0 {
			if j.execCount >= j.execLimit {
				log.Printf("Stopped job %v on epoch %v for reaching execution limit %v\n", j.name, j.execCount, j.execLimit)
				j.stopped = true
				s.jobWaitGroup.Done()
				return
			}
		}
		time.Sleep(time.Until(j.nextTime))
		startTime := time.Now()
		log.Printf("Starting epoch %v for job %v at time %v\n", j.execCount, j.name, startTime)
		j.job()
		endTime := time.Now()
		log.Printf("Ended epoch %v for job %v at time %v\n", j.execCount, j.name, endTime)
		log.Printf("Execution time for job %v on epoch %v was %v\n", j.name, j.execCount, endTime.Sub(startTime))
		j.execCount++
		j.updateJob()
	}
	log.Printf("Stopped job %v on epoch %v as main routine ended\n", j.name, j.execCount)
	j.stopped = true
	s.jobWaitGroup.Done()
}

func (s *Schedule) startJobRoutines() {
	for index := range s.jobQueue {
		s.jobWaitGroup.Add(1)
		go s.startJob(s.jobQueue[index])
	}
	s.jobWaitGroup.Wait()
}

func (s *Schedule) sortSchedule() {
	sort.Slice(s.jobQueue, func(i, j int) bool {
		return s.jobQueue[i].nextTime.Before(s.jobQueue[j].nextTime)
	})
}

func (s *Schedule) addJobToSchedule(job *Job) {
	job.createdAt = time.Now()
	job.execCount = 0
	job.id = job.computeHash()
	log.Printf("Adding job %v with ID %v to schedule at time %v\n", job.name, job.id, job.createdAt)
	s.jobQueue = append(s.jobQueue, job)
	s.sortSchedule()
}

func (s *Schedule) openLog() {
	f, err := os.OpenFile(time.Now().Format(time.RFC3339)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening output log file: %v", err)
	}
	log.SetOutput(f)
	log.Printf("Opened output log at time %v\n", time.Now())
}

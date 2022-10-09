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
	j.updateJob()
	for j.nextTime.After(time.Now()) {
		// This is because I'm yet unsure of Go's short circuit behavior
		if j.execLimit != 0 {
			if j.execCount >= j.execLimit {
				return
			}
		}
		time.Sleep(time.Until(j.nextTime))
		log.Print("=============================================================")
		log.Printf("Executing job %v on epoch %v:\t", j.name, j.execCount)
		j.job()
		log.Print("=============================================================")
		j.execCount++
		j.updateJob()
	}
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
	s.jobQueue = append(s.jobQueue, job)
	s.sortSchedule()
}

func (s *Schedule) openLog() {
	f, err := os.OpenFile(time.Now().Format(time.RFC3339)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(f)
}

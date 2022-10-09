package ipcron

import (
	"bytes"
	"log"
	"time"

	"crypto/md5"
	"encoding/gob"
	"encoding/hex"

	"github.com/aptible/supercronic/cronexpr"
)

type Job struct {
	id                 string
	name               string
	scheduleExpression string
	nextTime           time.Time
	interval           time.Duration
	job                func()
	execCount          int
	createdAt          time.Time
	execLimit          int
	stopped            bool
	maxRunTime         time.Duration
}

func (j *Job) computeHash() string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode([]int64{j.createdAt.UnixMicro()})
	hash := md5.Sum(b.Bytes())
	result := hex.EncodeToString(hash[:])
	log.Printf("Computed has for job %v as %v\n", j.name, result)
	return result
}

func (j *Job) updateJob() {
	log.Printf("Updating job %v at time %v\n", j.name, time.Now())
	if len(j.scheduleExpression) > 0 {
		j.nextTime = cronexpr.MustParse(j.scheduleExpression).Next(time.Now())
	} else {
		j.nextTime = time.Now().Add(j.interval)
	}
	log.Printf("Job %v next execution time at %v\n", j.name, j.nextTime)
}

func (j *Job) SetExecutionLimit(limit int) {
	log.Printf("Setting execution limit on job %v to %v runs\n", j.name, limit)
	j.execLimit = limit
}

func (j *Job) GetId() string {
	log.Printf("Retrieving ID for job %v as %v\n", j.name, j.id)
	return j.id
}

func (j *Job) GetExecutionCount() int {
	log.Printf("Retrieving execution count for job %v as %v\n", j.name, j.execCount)
	return j.execCount
}

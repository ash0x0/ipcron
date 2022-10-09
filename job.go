package ipcron

import (
	"bytes"
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
}

func (j *Job) computeHash() string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode([]int64{j.createdAt.UnixMicro()})
	hash := md5.Sum(b.Bytes())
	return hex.EncodeToString(hash[:])
}

func (j *Job) updateJob() {
	if len(j.scheduleExpression) > 0 {
		j.nextTime = cronexpr.MustParse(j.scheduleExpression).Next(time.Now())
	} else {
		j.nextTime = time.Now().Add(j.interval)
	}
}

func (j *Job) SetExecutionLimit(limit int) {
	j.execLimit = limit
}

func (j *Job) GetId() string {
	return j.id
}

func (j *Job) GetExecutionCount() int {
	return j.execCount
}

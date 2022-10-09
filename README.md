# IPcron: In-Process Cron for Go

A simpel and no-frills task scheduler for Go. Part interview task, part Go learning project, and a totally fun and nice experience for myself. Somewhat useless for everyone else ¯\_(ツ)_/¯. 

If you find yourself here and haven't been given a link, this is not the cron droid you're looking for.

# Getting Started

## tl;dr
In terminal:
```bash
go get github.com/ash0x0/ipcron
```

In your code:
```GO
schedule := ipcron.NewSchedule(true)    // get new scheduler process
job, err := schedule.ScheduleJobWithInterval(time.ParseDuration('10s'), example, "exampleJob")  // add a job with a simple time interval
job, _ := schedule.ScheduleWithCronSyntax("* * * * * * *", example, "cronExample")  // add a job with cron syntax
schedule.Start()    // start all added jobs, nothing will happen without this
schedule.Stop() // stop all jobs
```
---

There's a ready example in the `example` folder, check it out to see the gist of how this is supposed to be used.

In essence, this project allows you to schedule arbitrary functions to execute at preset times or intervals.

In either case, you need to create a scheduler to handle things for you first:
```Go
schedule := ipcron.NewSchedule(false)   // the false here means that logging won't be redirected to fil
schedule := ipcron.NewSchedule(true)   // this will create a new logfile and direct all logs there
```

Then you can start adding your functions to the scheduler at will. Imagine you have a simple function:
```Go
func example() {
    fmt.Println("The answer to life, the universe, and everything")
}
```

You have two options to add this to the scheduler. 

## Add With Interval
```Go
interval := time.ParseDuration('10s')
job, err := schedule.ScheduleJobWithInterval(interval, example, "exampleJob")
```

As you can see, this takes in a `time.Duration` type as the interval. The signature for this is as follows:

```Go
func (s *Schedule) ScheduleJobWithInterval(timeInterval time.Duration, job func(), jobName string) (*Job, error)
```

It will add the function such that there is duration `timeInterval` between the end of every execution and the beginning of the next. This is a bit more intuitive option than the usual cron "at this moment in time".

## Add With Cron Syntax
```Go
job, _ := schedule.ScheduleWithCronSyntax("* * * * * * *", example, "cronExample")
```

This allows adding with cron expression syntax, which is useful for generated and structured schedules, though definitely not as 
intuitive.

The signature for this is:
```Go
func (s *Schedule) ScheduleWithCronSyntax(scheduleExpression string, job func(), jobName string) (*Job, error)
```

When adding with cron syntax, the scheduling behavior follows the same expected cron behavior. Instead of scheduling jobs at difinite intervals, it executes the job at preset moments in time.

Whereas running with interval would say ***"execute this every 10 seconds"***, cron would say ***"execute this at multiples of 10 seconds on the clock"**

## Set Execution Limit

Any jobs you create will run forever, which is a very long time. To make things make a little more sense, you can set an execution limit count which will limit the job to a certain number of runs
```Go
job.SetExecutionLimit(10)
```
This needs to be run before starting the scheduler. I have no idea what will happen if it's set during the schedule run (it'll likely not do anything) and it's not supported.

## Start Schedule
After you add all your processes, you need to start the scheduler. No jobs will run without this.
```Go
scheduler.Start()
```

## Stop Jobs
If you wish to stop it at any time, just call `Stop()`
```Go
scheduler.Stop()
```

# Design Decisions

- Go was chosen because an "in-process cron scheduler" just sounds like something Go would be good at because concurrency. Also, I just wanted to learn some more Go.

- Cron scheduling syntax is made available through `supercronic/cronexpr/` package. I don't see it as a core function of scheduling, which is why I didn't implement it myself. Though cron syntax is useful for computer generated schedules and copy-paste, it's quite anti-intuitive. Creating a parser for it is also relatively straight-forward and trivial so I personally didn't care to implement it myself.

- **Jobs don't persist past process** - Obviously, given this is "in-process", it ends with the process. Although making it persist past the process would be relatively easy, it just breaks away from the concept so no attempt was ever made to create that.

- **Every job has a goroutine** - When the scheduler is started, each scheduled jobs gets its own goroutine where it lives for the remainder of the process. The job goroutines sleep when a job is waiting for the time to be run. The job goroutine should ideally execute and end when the job is due, not stick around sleeping. The reason for this choice is to make usage more predictable and the implementation simpler. It places control with the job goroutines and makes it so the scheduler goroutine doesn't have to act as a giant orchestrator, which would complicate its implementation.

- **Job Hash ID** - Jobs have an MD5 hash ID. I wanted jobs to have a unique identified because I was looking to implement a per-job `Stop()` with the ID, which ended up not making it in :(. Why not an integer? Because hashes are guaranteed to be unique while integer IDs are really up to the implementation, and I didn't wanna have to worry about making them unique.

# Limitations

- **Limited function signatrue** - This accepts functions with no input and output. `func(number int) string` will give you an error. Though an implementation with Go channels can make this possible, this is not the one opted for here. Therefore you can only execute basic functions with this implementation.

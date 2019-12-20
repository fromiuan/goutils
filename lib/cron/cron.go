// This library implements a cron spec parser and runner.  See the README for
// more details.
package cron

import (
	"fmt"
	"log"
	"runtime"
)

// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may be started, stopped, and the entries may
// be inspected while running.
type Cron struct {
	stop     chan struct{}
	add      chan *Entry
	running  bool
	ErrorLog *log.Logger
}

// Job is an interface for submitted cron jobs.
type Job interface {
	Run()
}

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
	// The Job to run.
	Job Job
}

// New returns a new Cron job runner.
func NewCron(size int) *Cron {
	return &Cron{
		add:      make(chan *Entry, size),
		stop:     make(chan struct{}),
		running:  false,
		ErrorLog: nil,
	}

}

// AddJob adds a Job to the Cron to be run on the given schedule.
func (c *Cron) AddJob(cmd Job) error {
	entry := &Entry{
		Job: cmd,
	}
	select {
	case c.add <- entry:
		return nil
	default:
		return fmt.Errorf("job cron queue is full")
	}
}

// Start the cron scheduler in its own go-routine.
func (c *Cron) Start() {
	c.running = true
	go c.run()
}

func (c *Cron) runWithRecovery(entry *Entry) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.logf("cron: panic running job: %v\n%s", r, buf)
		}
	}()
	entry.Job.Run()
}

// Run the scheduler.. this is private just due to the need to synchronize
// access to the 'running' state variable.
func (c *Cron) run() {
	for {
		select {
		case entry, ok := <-c.add:
			if ok {
				c.runWithRecovery(entry)
			} else {
				c.logf("Cron error")
			}
		case <-c.stop:
			c.logf("Cron stop")
			return
		}
	}
}

// Logs an error to stderr or to the configured error log
func (c *Cron) logf(format string, args ...interface{}) {
	if c.ErrorLog != nil {
		c.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing.
func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

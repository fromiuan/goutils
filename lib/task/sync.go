package task

import (
	"log"
	"strings"
	"sync"
	"time"
)

// SyncTask task struct
type SyncTask struct {
	Taskname string
	Spec     *Schedule
	SpecStr  string
	DoFunc   TaskFunc
	Prev     time.Time
	Next     time.Time
	Errlist  []*taskerr // like errtime:errinfo
	ErrLimit int        // max length for the errlist, 0 stand for no limit
	running  bool
	lock     sync.Mutex
}

// NewSyncTask add new task with name, time and func
func NewSyncTask(tname string, spec string, f TaskFunc) *SyncTask {

	task := &SyncTask{
		Taskname: tname,
		DoFunc:   f,
		ErrLimit: 100,
		SpecStr:  spec,
		running:  false,
	}
	task.SetCron(spec)
	return task
}

// GetSpec get spec string
func (t *SyncTask) GetSpec() string {
	return t.SpecStr
}

// GetStatus get current task status
func (t *SyncTask) GetStatus() string {
	var str string
	for _, v := range t.Errlist {
		str += v.t.String() + ":" + v.errinfo + "<br>"
	}
	return str
}

// Run run all tasks
func (t *SyncTask) Run() error {
	t.lock.Lock()
	if t.running {
		t.lock.Unlock()
		return nil
	}
	t.running = true
	t.lock.Unlock()
	err := t.DoFunc()
	if err != nil {
		if t.ErrLimit > 0 && t.ErrLimit > len(t.Errlist) {
			t.Errlist = append(t.Errlist, &taskerr{t: t.Next, errinfo: err.Error()})
		}
	}
	t.lock.Lock()
	t.running = false
	t.lock.Unlock()
	return err
}

// SetNext set next time for this task
func (t *SyncTask) SetNext(now time.Time) {
	t.Next = t.Spec.Next(now)
}

// GetNext get the next call time of this task
func (t *SyncTask) GetNext() time.Time {
	return t.Next
}

// SetPrev set prev time of this task
func (t *SyncTask) SetPrev(now time.Time) {
	t.Prev = now
}

// GetPrev get prev time of this task
func (t *SyncTask) GetPrev() time.Time {
	return t.Prev
}

// six columns mean：
//       second：0-59
//       minute：0-59
//       hour：1-23
//       day：1-31
//       month：1-12
//       week：0-6（0 means Sunday）

// SetCron some signals：
//       *： any time
//       ,：　 separate signal
//　　    －：duration
//       /n : do as n times of time duration
/////////////////////////////////////////////////////////
//	0/30 * * * * *                        every 30s
//	0 43 21 * * *                         21:43
//	0 15 05 * * * 　　                     05:15
//	0 0 17 * * *                          17:00
//	0 0 17 * * 1                           17:00 in every Monday
//	0 0,10 17 * * 0,2,3                   17:00 and 17:10 in every Sunday, Tuesday and Wednesday
//	0 0-10 17 1 * *                       17:00 to 17:10 in 1 min duration each time on the first day of month
//	0 0 0 1,15 * 1                        0:00 on the 1st day and 15th day of month
//	0 42 4 1 * * 　 　                     4:42 on the 1st day of month
//	0 0 21 * * 1-6　　                     21:00 from Monday to Saturday
//	0 0,10,20,30,40,50 * * * *　           every 10 min duration
//	0 */10 * * * * 　　　　　　              every 10 min duration
//	0 * 1 * * *　　　　　　　　               1:00 to 1:59 in 1 min duration each time
//	0 0 1 * * *　　　　　　　　               1:00
//	0 0 */1 * * *　　　　　　　               0 min of hour in 1 hour duration
//	0 0 * * * *　　　　　　　　               0 min of hour in 1 hour duration
//	0 2 8-20/3 * * *　　　　　　             8:02, 11:02, 14:02, 17:02, 20:02
//	0 30 5 1,15 * *　　　　　　              5:30 on the 1st day and 15th day of month
func (t *SyncTask) SetCron(spec string) {
	t.Spec = t.parse(spec)
}

func (t *SyncTask) parse(spec string) *Schedule {
	if len(spec) > 0 && spec[0] == '@' {
		return t.parseSpec(spec)
	}
	// Split on whitespace.  We require 5 or 6 fields.
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	fields := strings.Fields(spec)
	if len(fields) != 5 && len(fields) != 6 {
		log.Panicf("Expected 5 or 6 fields, found %d: %s", len(fields), spec)
	}

	// If a sixth field is not provided (DayOfWeek), then it is equivalent to star.
	if len(fields) == 5 {
		fields = append(fields, "*")
	}

	schedule := &Schedule{
		Second: getField(fields[0], seconds),
		Minute: getField(fields[1], minutes),
		Hour:   getField(fields[2], hours),
		Day:    getField(fields[3], days),
		Month:  getField(fields[4], months),
		Week:   getField(fields[5], weeks),
	}

	return schedule
}

func (t *SyncTask) parseSpec(spec string) *Schedule {
	switch spec {
	case "@yearly", "@annually":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  1 << months.min,
			Week:   all(weeks),
		}

	case "@monthly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@weekly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   1 << weeks.min,
		}

	case "@daily", "@midnight":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@hourly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   all(hours),
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}
	}
	log.Panicf("Unrecognized descriptor: %s", spec)
	return nil
}

// Name return task name
func (t *SyncTask) Name() string {
	return t.Taskname
}

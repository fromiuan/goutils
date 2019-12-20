package cron

import (
	"log"
)

type Task interface {
	Run() error
}

type TaskJob struct {
	Task Task
}

type Worker struct {
	WorkerPool     chan chan TaskJob
	TaskJobChannel chan TaskJob
	quit           chan bool
}

type Dispatcher struct {
	WorkerPool   chan chan TaskJob
	TaskJobQueue chan TaskJob
	maxWorkers   int
}

func NewWorker(workerPool chan chan TaskJob) Worker {
	return Worker{
		WorkerPool:     workerPool,
		TaskJobChannel: make(chan TaskJob),
		quit:           make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.TaskJobChannel
			select {
			case taskJob := <-w.TaskJobChannel:
				err := taskJob.Task.Run()
				if err != nil {
					log.Println(err)
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func NewDispatcher(maxWorker, maxQueue int) *Dispatcher {
	pool := make(chan chan TaskJob, maxWorker)
	return &Dispatcher{
		WorkerPool:   pool,
		maxWorkers:   maxWorker,
		TaskJobQueue: make(chan TaskJob, maxQueue),
	}
}

func (d *Dispatcher) AddJob(t Task) {
	job := TaskJob{
		Task: t,
	}
	d.TaskJobQueue <- job
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		w := NewWorker(d.WorkerPool)
		w.Start()
	}
	d.dispatcher()
}

func (d *Dispatcher) dispatcher() {
	go func() {
		for {
			select {
			case taskJob := <-d.TaskJobQueue:
				go func(taskJob TaskJob) {
					taskJobChan := <-d.WorkerPool
					taskJobChan <- taskJob
				}(taskJob)
			}
		}
	}()

}

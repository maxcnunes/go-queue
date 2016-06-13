package main

import (
	"log"
	"sync"
)

type Queue struct {
	Jobs   chan interface{} // Jobs is the channel that stores the items to be processed
	Errors chan error       // Errors is the channel that handle the errors
	task   QueueTaskRunner
	wg     sync.WaitGroup
}

// QueueTaskRunner expect the implementation to process the job
type QueueTaskRunner interface {
	Run(job interface{}) error
}

// NewQueue creates a queue
func NewQueue(concurrent int, task QueueTaskRunner) *Queue {
	wg := sync.WaitGroup{}
	wg.Add(concurrent)

	return &Queue{
		Jobs:   make(chan interface{}, concurrent),
		Errors: make(chan error),
		task:   task,
		wg:     wg,
	}
}

// Run executes the queue with n workers based in the concurrent number
func (q *Queue) Run() *Queue {
	for w := cap(q.Jobs); w > 0; w-- {
		go q.worker()
	}

	return q
}

// Wait all jobs be processed
func (q *Queue) Drain() error {
	q.wg.Wait()
	if len(q.Errors) > 0 {
		return <-q.Errors
	}
	return nil
}

func (q *Queue) worker() {
	for {
		select {
		case job, ok := <-q.Jobs:
			if !ok {
				q.wg.Done()
				return
			}

			log.Printf("Processing job %#v", job)
			if err := q.task.Run(job); err != nil {
				q.wg.Done()
				q.Errors <- err
				return
			}
		}
	}
}

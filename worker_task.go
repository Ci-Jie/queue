package queue

import (
	"context"
	"errors"
)

var _ Worker = (*taskWorker)(nil)

// just for unit testing, don't use it.
type taskWorker struct {
	messages chan QueuedMessage
}

func (w *taskWorker) BeforeRun() error { return nil }
func (w *taskWorker) AfterRun() error  { return nil }
func (w *taskWorker) Run() error {
	for msg := range w.messages {
		if v, ok := msg.(Job); ok {
			if v.Task != nil {
				_ = v.Task(context.Background())
			}
		}
	}
	return nil
}

func (w *taskWorker) Shutdown() error {
	close(w.messages)
	return nil
}

func (w *taskWorker) Queue(job QueuedMessage) error {
	select {
	case w.messages <- job:
		return nil
	default:
		return errors.New("max capacity reached")
	}
}
func (w *taskWorker) Capacity() int       { return cap(w.messages) }
func (w *taskWorker) Usage() int          { return len(w.messages) }
func (w *taskWorker) BusyWorkers() uint64 { return uint64(0) }

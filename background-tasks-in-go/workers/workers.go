package workers

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

type worker struct {
	workchan    chan workType
	workerCount int
	buffer      int
	wg          *sync.WaitGroup
	cancelFunc  context.CancelFunc
}

type WorkerIface interface {
	Start(pctx context.Context)
	Stop()
	QueueTask(task string, workDuration time.Duration) error
}

func New(workerCount, buffer int) WorkerIface {
	w := worker{
		workchan:    make(chan workType, buffer),
		workerCount: workerCount,
		buffer:      buffer,
		wg:          new(sync.WaitGroup),
	}

	return &w
}

func (w *worker) Start(pctx context.Context) {
	ctx, cancelFunc := context.WithCancel(pctx)
	w.cancelFunc = cancelFunc

	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.spawnWorkers(ctx)
	}
}

func (w *worker) Stop() {
	log.Info("stop workers")
	close(w.workchan)
	w.cancelFunc()
	w.wg.Wait()
	log.Info("all workers exited!")
}

func (w *worker) QueueTask(task string, workDuration time.Duration) error {
	if len(w.workchan) >= w.buffer {
		return ErrWorkerBusy
	}

	w.workchan <- workType{TaskID: task, WorkDuration: workDuration}

	return nil
}

func (w *worker) spawnWorkers(ctx context.Context) {
	defer w.wg.Done()

	for work := range w.workchan {
		select {
		case <-ctx.Done():
			return
		default:
			w.doWork(ctx, work.TaskID, work.WorkDuration)
		}
	}
}

func (w *worker) doWork(ctx context.Context, task string, workDuration time.Duration) {
	log.WithField("task", task).Info("do some work now...")
	sleepContext(ctx, workDuration)
	log.WithField("task", task).Info("work completed!")
}

func sleepContext(ctx context.Context, sleep time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(sleep):
	}
}

type workType struct {
	TaskID       string
	WorkDuration time.Duration
}

var (
	ErrWorkerBusy = errors.New("workers are busy, try again later")
)

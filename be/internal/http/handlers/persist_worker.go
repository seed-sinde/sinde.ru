package handlers

import (
	"context"
	"log"
	"sync"
	"time"
)

type persistTask struct {
	name string
	fn   func(ctx context.Context)
}

var (
	persistOnce  sync.Once
	persistQueue chan persistTask
)

const (
	persistQueueSize   = 256
	persistWorkerCount = 4
	persistTaskTimeout = 5 * time.Second
)

func initPersistWorkers() {
	persistQueue = make(chan persistTask, persistQueueSize)
	for i := 0; i < persistWorkerCount; i++ {
		go func(workerID int) {
			for task := range persistQueue {
				ctx, cancel := context.WithTimeout(context.Background(), persistTaskTimeout)
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[воркер персиста %d] panic в задаче %s: %v", workerID, task.name, r)
						}
					}()
					task.fn(ctx)
				}()
				cancel()
			}
		}(i)
	}
}
func enqueuePersistTask(name string, fn func(ctx context.Context)) {
	persistOnce.Do(initPersistWorkers)
	select {
	case persistQueue <- persistTask{name: name, fn: fn}:
	default:
		log.Printf("[персист] очередь заполнена, задача %s отбрасывается", name)
	}
}

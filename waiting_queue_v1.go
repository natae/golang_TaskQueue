package waitingqueue

import (
	"sync"
)

type WaitingQueueV1 struct {
	sync.RWMutex
	taskQueue []Task
}

func (wQueue *WaitingQueueV1) Start() {
	go wQueue.doTask_goroutine()
}

func (wQueue *WaitingQueueV1) Destroy() {

}

func (wQueue *WaitingQueueV1) doTask_goroutine() {
	for {
		if len(wQueue.taskQueue) <= 0 {
			continue
		}

		task := wQueue.taskQueue[0]
		result := task.Func(task.Params...)
		task.ReturnChan <- result

		wQueue.Lock()
		wQueue.taskQueue = wQueue.taskQueue[1:]
		wQueue.Unlock()
	}
}

func (wQueue *WaitingQueueV1) RequestTask(task Task) int {
	wQueue.Lock()
	defer wQueue.Unlock()
	if taskIndex := wQueue.existTask(task.Id); taskIndex < 0 {
		wQueue.taskQueue = append(wQueue.taskQueue, task)

		return len(wQueue.taskQueue)
	} else {
		return taskIndex
	}
}

func (wQueue *WaitingQueueV1) existTask(id string) int {
	index := 0
	for _, task := range wQueue.taskQueue {
		if task.Id == id {
			return index
		}
		index++
	}

	return -1
}

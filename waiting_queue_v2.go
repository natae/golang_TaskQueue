package waitingqueue

type WaitingQueueV2 struct {
	checkTaskChan chan Task
	doTaskChan    chan Task
	endTaskChan   chan struct{}
	taskIdQueue   []string
	resultChan    chan int
}

func (wQueue *WaitingQueueV2) Start() {
	go wQueue.checkTask_goroutine()
	go wQueue.doTask_goroutine()
}

func (wQueue *WaitingQueueV2) Destroy() {
	close(wQueue.checkTaskChan)
	close(wQueue.doTaskChan)
	close(wQueue.endTaskChan)
	close(wQueue.resultChan)
}

func (wQueue *WaitingQueueV2) doTask_goroutine() {
	for {
		task := <-wQueue.doTaskChan

		result := task.Func(task.Params...)
		task.ReturnChan <- result

		wQueue.endTaskChan <- struct{}{}
	}
}

func (wQueue *WaitingQueueV2) checkTask_goroutine() {
	for {
		select {
		case task := <-wQueue.checkTaskChan:

			if index := wQueue.existTask(task.Id); index < 0 {
				// New task
				wQueue.taskIdQueue = append(wQueue.taskIdQueue, task.Id)

				wQueue.resultChan <- len(wQueue.taskIdQueue)

				// Throw task
				wQueue.doTaskChan <- task
			} else {
				// Existing task
				wQueue.resultChan <- index
			}
		case <-wQueue.endTaskChan:
			wQueue.taskIdQueue = wQueue.taskIdQueue[1:]
		}
	}
}

func (wQueue *WaitingQueueV2) existTask(targetTaskId string) int {
	index := 0
	for _, taskId := range wQueue.taskIdQueue {
		if taskId == targetTaskId {
			return index
		}
		index++
	}

	return -1
}

func (wQueue *WaitingQueueV2) RequestTask(task Task) int {

	wQueue.checkTaskChan <- task

	return <-wQueue.resultChan
}

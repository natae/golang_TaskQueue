package waitingqueue

import "sync"

type IWaitingQueue interface {
	Start()
	Destroy()
	RequestTask(task Task) int
}

const (
	DefaultQueueLength = 1024
)

type Config struct {
	QueueLength  int
	IsOldVersion bool
}

func New(config ...Config) IWaitingQueue {

	targetConfig := Config{}

	if len(config) > 0 {
		targetConfig = config[0]
	}

	if targetConfig.QueueLength <= 0 {
		targetConfig.QueueLength = DefaultQueueLength
	}

	if targetConfig.IsOldVersion {
		waitingQueueV1 := &WaitingQueueV1{}

		waitingQueueV1.RWMutex = sync.RWMutex{}
		waitingQueueV1.taskQueue = make([]Task, 0, targetConfig.QueueLength)

		return waitingQueueV1
	} else {
		waitingQueueV2 := &WaitingQueueV2{}

		waitingQueueV2.checkTaskChan = make(chan Task)
		waitingQueueV2.doTaskChan = make(chan Task, targetConfig.QueueLength)
		waitingQueueV2.endTaskChan = make(chan struct{}, targetConfig.QueueLength)
		waitingQueueV2.taskIdQueue = make([]string, 0, targetConfig.QueueLength)
		waitingQueueV2.resultChan = make(chan int)

		return waitingQueueV2
	}
}

package waitingqueue

type TaskFunc func(params ...interface{}) interface{}

type Task struct {
	Id         string
	Func       TaskFunc
	Params     []interface{}
	ReturnChan chan interface{}
}

func NewTask(id string, taskFunc TaskFunc, params []interface{}) Task {
	return Task{
		Id:         id,
		Func:       taskFunc,
		Params:     params,
		ReturnChan: make(chan interface{}),
	}
}

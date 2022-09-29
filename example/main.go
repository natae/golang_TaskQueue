package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/natae/waitingqueue"
)

func main() {
	waitingQueue := waitingqueue.New()
	waitingQueue.Start()

	waitGroup := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		waitGroup.Add(1)
		go func(index int) {
			task := waitingqueue.NewTask(strconv.Itoa(index), pow, []interface{}{index, 2})

			for {
				select {
				case ret := <-task.ReturnChan:
					fmt.Printf("[%d] Return value: %d\n", index, ret)
					waitGroup.Done()
					return
				default:
					remains := waitingQueue.RequestTask(task)
					fmt.Printf("[%d] Wait order: %d\n", index, remains)
				}
				time.Sleep(time.Millisecond * 500)
			}

		}(i)
	}

	waitGroup.Wait()

	waitingQueue.Destroy()
}

func pow(params ...interface{}) interface{} {
	if len(params) < 2 {
		return 0
	}

	num := params[0].(int)
	n := params[1].(int)
	result := 1
	for i := 0; i < n; i++ {
		result *= num
		time.Sleep(time.Second * 1)
	}
	return result
}

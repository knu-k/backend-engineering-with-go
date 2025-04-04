package task

import (
	"fmt"
	"sync"
	"time"
)

// ###### Synchronous Task
func syncTask(id int) {
	fmt.Printf("Task %d is running synchronously\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Task %d completed\n", id)
}

func TestSyncTaskDuration(repeat int) float64 {
	start := time.Now()
	for i := 1; i <= repeat; i++ {
		syncTask(i)
	}
	end := time.Since(start)
	return end.Seconds()
}

// ###### Asynchronous Task
func asyncTask(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Task %d is running asynchronously\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Task %d completed\n", id)
}
func TestAsyncTaskDuration(repeat int) float64 {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 1; i <= repeat; i++ {
		wg.Add(1)
		go asyncTask(i, &wg)
	}
	wg.Wait()
	end := time.Since(start)
	return end.Seconds()
}
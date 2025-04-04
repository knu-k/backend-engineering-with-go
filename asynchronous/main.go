package main

import (
	"fmt"
	"runtime"

	"github.com/polyglot-k/asynchronous/task"
)


func main() {
	runtime.GOMAXPROCS(1) // OS 스레드 1개로 제한
	syncDuration := task.TestSyncTaskDuration(5)
	asyncDuration := task.TestAsyncTaskDuration(5)
	fmt.Printf("Synchronous tasks took %.2f seconds\n", syncDuration)
	fmt.Printf("Asynchronous tasks took %.2f seconds\n", asyncDuration)
}
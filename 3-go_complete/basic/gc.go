package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"runtime/trace"
	"sync"
	"time"
)

var arr [1000]int                 //数组长度固定，通常在栈上
var slc []int = make([]int, 1000) //动态创建的结构体、切片、map在堆上。引用类型的全局变量内存分配在堆上，值类型的全局变量分配在栈上

func stack_heap() {
	var brr [1000]int

	var crr [1310720]int
	var drr [1310721]int //数组超过10M就会分配到堆上(moved to heap)

	err := make([]int, 8192) //函数的入参、出参、局部变量一般在栈上
	frr := make([]int, 8193) //切片超过64K就会分配到堆上(escapes to heap)

	_ = arr
	_ = brr
	_ = crr
	_ = drr
	_ = err
	_ = frr
	_ = slc
}

const (
	NumWorkers    = 4     // Number of workers.
	NumTasks      = 500   // Number of tasks.
	MemoryIntense = 10000 // Size of memory-intensive task (number of elements).
)

// 所有的垃圾回收都是针对堆的
func gc() {
	// Write to the trace file.
	f, _ := os.Create("data/trace.out")
	trace.Start(f)
	defer trace.Stop()

	// Set the target percentage for the garbage collector. Default is 100%.
	debug.SetGCPercent(100)

	// Task queue and result queue.
	taskQueue := make(chan int, NumTasks)
	resultQueue := make(chan int, NumTasks)

	// Start workers.
	var wg sync.WaitGroup
	wg.Add(NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		go worker(taskQueue, resultQueue, &wg)
	}

	// Send tasks to the queue.
	for i := 0; i < NumTasks; i++ {
		taskQueue <- i
	}
	close(taskQueue)

	// Retrieve results from the queue.
	go func() {
		wg.Wait()
		close(resultQueue)
	}()

	// Process the results.
	for result := range resultQueue {
		fmt.Println("Result:", result)
	}

	fmt.Println("Done!")
}

// Worker function.
func worker(tasks <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		result := performMemoryIntensiveTask(task)
		results <- result
	}
}

// performMemoryIntensiveTask is a memory-intensive function.
func performMemoryIntensiveTask(task int) int {
	// Create a large-sized slice.
	data := make([]int, MemoryIntense)
	for i := 0; i < MemoryIntense; i++ {
		data[i] = i + task
	}

	// Latency imitation.
	time.Sleep(10 * time.Millisecond)

	// Calculate the result.
	result := 0
	for _, value := range data {
		result += value
	}
	return result
}

func main16() {
	// stack_heap() // go build -gcflags=-m ./basic/gc.go
	gc() // 程序运行完之后生成一个文件data/trace.out, 然后执行 go tool trace data/trace.out
}

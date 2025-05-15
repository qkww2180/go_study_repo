package d_concurrence

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	READ_FILE_ROUTINE_COUNT = 15
	PROCESS_ROUTINE_COUNT   = 5
)

var (
	sum int64

	fileList   = make(chan string, 100)
	lineBuffer = make(chan string, 1000)

	walkWg    sync.WaitGroup
	readWg    sync.WaitGroup
	processWg sync.WaitGroup
)

// 递归遍历dir目录，把文件全都放入fileList
func walkDir(dir string) {
	defer walkWg.Done()
	filepath.Walk(dir, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.Mode().IsRegular() {
			fileList <- subPath //如果下游消费得慢，这里可能会阻塞
		}
		return nil
	})
}

// 读完一个文件，读下一个文件，纯IO操作，可以多开几个协程
func readFile() {
	defer readWg.Done()
	for {
		if infile, ok := <-fileList; ok {
			fin, err := os.Open(infile)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer fin.Close()
			reader := bufio.NewReader(fin)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						if len(line) > 0 {
							lineBuffer <- strings.TrimSpace(line) //把前、后的空格（包括换行符）删除掉
						}
						break
					} else {
						fmt.Println(err)
					}
				} else {
					lineBuffer <- strings.TrimSpace(line)
				}
			}
		} else {
			break
		}
	}
}

func processLine() {
	defer processWg.Done()
	for {
		if line, ok := <-lineBuffer; ok {
			if i, err := strconv.Atoi(line); err != nil {
				fmt.Printf("%s not number\n", line)
			} else {
				atomic.AddInt64(&sum, int64(i))
			}
		} else {
			break
		}
	}
}

// 并行处理海量文件
func DealMassFile(dir string) {

	go func() {
		tk := time.NewTicker(time.Second)
		defer tk.Stop()
		for {
			<-tk.C
			fmt.Printf("堆积了%d个文件未处理，堆积了%d行内容未处理\n", len(fileList), len(lineBuffer)) //如果堆积得多，就需要加大下游的并发度，但如果CPU利用率已经很高了，可以减小上游的并发度
		}
	}()

	walkWg.Add(1)
	readWg.Add(READ_FILE_ROUTINE_COUNT)
	processWg.Add(PROCESS_ROUTINE_COUNT)

	go walkDir(dir)
	for i := 0; i < READ_FILE_ROUTINE_COUNT; i++ {
		go readFile()
	}
	for i := 0; i < PROCESS_ROUTINE_COUNT; i++ {
		go processLine()
	}

	walkWg.Wait()
	close(fileList)
	readWg.Wait()
	close(lineBuffer)
	processWg.Wait()

	fmt.Printf("sum=%d\n", sum)
}

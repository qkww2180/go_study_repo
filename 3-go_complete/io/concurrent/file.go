/*
并行读写文件
*/
package concurrent

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	log = "Rewind's personal plan for backing up unlimited repositories and cloud storage. Includes 14-day free trial.\n"
)

// 同一个文件支持在多个协程里同时被打开，并发写入
func MultiFilterWriter(outFile string) {
	const C = 10
	wg := sync.WaitGroup{}
	wg.Add(C)
	for i := 0; i < C; i++ {
		go func() {
			defer wg.Done()
			fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644) //由于要被打开多次，所以就使用O_APPEND，而非O_TUNC
			if err != nil {
				panic(err)
			}
			defer fout.Close()
			for j := 0; j < 1000; j++ {
				fout.WriteString(log)
			}
		}()
	}
	wg.Wait()
}

// os.File本身支持并发调用
func OneFilterWriter(fout *os.File) {
	const C = 10
	wg := sync.WaitGroup{}
	wg.Add(C)
	for i := 0; i < C; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				fout.WriteString(log)
			}
		}()
	}
	wg.Wait()
}

// 把一个文件均等地分成互不相交的n份，且每份都是从新的一行开始。返回每一份的开始位置
func GetConcurrentReadPosition(infile string, n int) (positions []int64, err error) {
	begin := time.Now()
	defer func() {
		fmt.Printf("GetConcurrentReadPosition use time %d ms, begin positions %v\n", time.Since(begin).Milliseconds(), positions)
	}()
	fin, err := os.Open(infile)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	stat, err := fin.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stat.Size() //文件总大小
	fmt.Printf("文件大小 %d M\n", fileSize/1024/1024)
	size := fileSize / int64(n)                 //每一份的大小
	dupPositions := make(map[int64]struct{}, n) //确保positions中不存在重复的元素
	positions = append(positions, 0)            //第一份肯定是从位置0开始
	dupPositions[0] = struct{}{}

	var i int64
	buffer := make([]byte, 1024) //一次读1K
	for i = 1; i < int64(n); i++ {
		blockBegin := size * i
		//先定位到大概的位置，再逐个字节地往后找，找到第一个\n为止
		fin.Seek(blockBegin, 0) //whence=0相对于文件的起始位置。使用Seek的时候不要用APPEND模式打开文件
		var foundReturn bool    //是否找到了换行符
		var pos int64           //第一个换行符的位置
	LOOP:
		for j := 0; j < 10; j++ { //一次读1K，我们认为10k之内一定能找到\n。如果真找不到就不往positions里添加元素
			n, err := fin.Read(buffer) //读出1k
			if err != nil {
				if err == io.EOF {
					break LOOP //读到了文件末尾都没有找到\n
				} else {
					return nil, err
				}
			}
			for i, ele := range buffer[:n] { //在这1k里面逐个字节地找\n
				if ele == '\n' {
					pos += int64(i)
					foundReturn = true
					break LOOP
				}
			}
			pos += int64(n)
		}
		if foundReturn {
			ele := blockBegin + pos + 1 //+1表示从\n的下一个位置开始
			if ele >= fileSize {
				break
			}
			if _, exists := dupPositions[ele]; exists {
				break
			} else {
				// fmt.Println(ele)
				dupPositions[ele] = struct{}{}
				positions = append(positions, ele)
			}
		} else {
			//10K之内没有找到\n，就不往positions里添加元素
		}
	}
	return positions, nil
}

// 处理文件，最多只读readTotal个字节
func processFile(fin *os.File, readTotal int64) {
	reader := bufio.NewReader(fin)
	var readCnt int64
	for {
		log, err := reader.ReadString('\n')
		if readTotal > 0 {
			readCnt += int64(len(log))
			if readCnt > readTotal {
				break
			}
		}
		if err != nil {
			if err == io.EOF {
				if len(log) > 0 {
					n := ExtractNumber(log)
					if n >= 0 {
						// sum += int32(n)
						atomic.AddInt32(&sum, int32(n)) //并发安全
					}
				}
			} else {
				fmt.Printf("读文件失败:%v", err)
			}
			break
		} else {
			log = strings.TrimRight(log, "\n")
			if len(log) > 0 {
				n := ExtractNumber(log)
				if n >= 0 {
					// sum += int32(n)
					atomic.AddInt32(&sum, int32(n)) //并发安全
				}
			}
		}
	}
}

// 把一个文件分成concurrent段，并行读取
func ConcurrentReadOneFile(infile string, n int, f func(*os.File, int64)) error {
	positions, err := GetConcurrentReadPosition(infile, n)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(positions))
	for i, position := range positions {
		currPosition := position
		var nextPositoin int64
		if i < len(positions)-1 {
			nextPositoin = positions[i+1]
		} else {
			nextPositoin = -1
		}
		go func(currPosition, nextPositoin int64) {
			defer wg.Done()
			fin, err := os.Open(infile)
			if err != nil {
				fmt.Printf("open file %s failed: %v", infile, err)
				return
			}
			defer fin.Close()
			fin.Seek(currPosition, 0)
			if nextPositoin > currPosition {
				f(fin, nextPositoin-currPosition)
			} else {
				f(fin, 0)
			}
		}(currPosition, nextPositoin)
	}
	wg.Wait()
	return nil
}

package main

import (
	"dqq/util"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	ARRAY_SIZE = 1e7 //1e8会超出gorutine对栈大小的使用限制
)

var (
	arr   = make([]byte, ARRAY_SIZE)
	index = [ARRAY_SIZE]int{}
)

func init() {
	for i := 0; i < len(index); i++ {
		index[i] = rand.Intn(len(arr))
	}
}

func writeRamSequentially() {
	t0 := time.Now()
	defer func() {
		fmt.Printf("顺序写内存%dms\n", time.Since(t0).Milliseconds())
	}()

	for i := 0; i < len(arr); i++ {
		arr[i] = byte(i)
	}
}

func readRamSequentially() {
	t0 := time.Now()
	defer func() {
		fmt.Printf("顺序读内存%dms\n", time.Since(t0).Milliseconds())
	}()

	for i := 0; i < len(arr); i++ {
		_ = arr[i]
	}
}

func readRamRandomly() {
	t0 := time.Now()
	defer func() {
		fmt.Printf("随机读内存%dms\n", time.Since(t0).Milliseconds())
	}()

	for _, i := range index {
		_ = arr[i]
	}
}

func writeRamRandomly() {
	t0 := time.Now()
	defer func() {
		fmt.Printf("随机写内存%dms\n", time.Since(t0).Milliseconds())
	}()

	for _, i := range index {
		arr[i] = byte(i)
	}
}

func writeDiskSequentially() {
	t0 := time.Now()
	defer func() {
		fmt.Printf("顺序写磁盘%dms\n", time.Since(t0).Milliseconds())
	}()

	fout, err := os.OpenFile(util.RootPath+"/data/arr.bin", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	fout.Write(arr)
}

func readDiskSequentially() {
	t0 := time.Now()
	defer func() {
		fmt.Printf("顺序读磁盘%dms\n", time.Since(t0).Milliseconds())
	}()

	fin, err := os.Open(util.RootPath + "/data/arr.bin")
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	fin.Read(arr)
}

func main() {
	writeRamSequentially()  //8ms
	readRamSequentially()   //5ms
	writeRamRandomly()      //65ms
	readRamRandomly()       //41ms
	writeDiskSequentially() //7ms
	readDiskSequentially()  //6ms
}

// go run .\io\sequential_io\

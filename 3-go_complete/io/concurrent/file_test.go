package concurrent

import (
	"dqq/io/util"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMultiFilterWriter(t *testing.T) {
	outFile := util.RootPath + "data/file_con_mw.log"
	os.Remove(outFile)
	begin := time.Now()
	MultiFilterWriter(outFile)
	fmt.Printf("MultiFilterWriter time %dms\n", time.Since(begin).Milliseconds()) //96ms
}

func TestOneFilterWriter(t *testing.T) {
	outFile := util.RootPath + "data/file_con_sw.log"
	os.Remove(outFile)
	begin := time.Now()
	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	OneFilterWriter(fout)
	fout.Close()
	fmt.Printf("OneFilterWriter time %dms\n", time.Since(begin).Milliseconds()) //39ms
}

func TestSequentialWriter(t *testing.T) {
	outFile := util.RootPath + "data/file_seq.log"
	os.Remove(outFile)
	begin := time.Now()
	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10000; i++ {
		fout.WriteString(log)
	}
	fout.Close()
	fmt.Printf("SequentialWriter time %dms\n", time.Since(begin).Milliseconds()) //25ms
}

func TestConcurrentReadOneFile(t *testing.T) {
	infile := "/data1/search_log/search.log.20230505"
	sum = 0
	begin := time.Now()
	const P = 4 //4个协程并行处理文件
	if err := ConcurrentReadOneFile(infile, P, processFile); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("并行解析文件sum=%d, time %dms\n", sum, time.Since(begin).Milliseconds())
}

func TestSequentialReadOneFile(t *testing.T) {
	infile := "/data1/search_log/search.log.20230505"
	sum = 0
	begin := time.Now()
	fin, err := os.Open(infile)
	if err != nil {
		fmt.Printf("open file %s failed: %v", infile, err)
		return
	}
	defer fin.Close()
	processFile(fin, 0)
	fmt.Printf("串行解析文件sum=%d, time %dms\n", sum, time.Since(begin).Milliseconds())
}

//go test -v .\io\concurrent\ -run=^TestMultiFilterWriter$ -count=1
//go test -v .\io\concurrent\ -run=^TestOneFilterWriter$ -count=1
//go test -v .\io\concurrent\ -run=^TestSequentialWriter$ -count=1
//go test -v .\io\concurrent\ -run=^TestConcurrentReadOneFile$ -count=1
//go test -v .\io\concurrent\ -run=^TestSequentialReadOneFile$ -count=1

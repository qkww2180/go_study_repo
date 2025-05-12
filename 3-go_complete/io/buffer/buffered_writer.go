package main

import (
	"dqq/util"
	"fmt"
	"os"
	"time"
)

var (
	text = "recall 74 by radic, recall 0 by milvus_short, recall 0 by milvus_long\n"
)

// 带缓冲的FileWriter
//
// Note: 不支持并发。golang自带的bufio.NewWriter支持并发
type BufferedFileWriter struct {
	buffer         []byte   //缓存的内容
	bufferEndIndex int      //buffer里有效内容的结束位置
	fout           *os.File //文件句柄
}

func NewWriter(fout *os.File, bufferSize int) *BufferedFileWriter {
	return &BufferedFileWriter{
		buffer:         make([]byte, bufferSize), //len=cap=bufferSize
		bufferEndIndex: 0,
		fout:           fout,
	}
}

// 向文件中写入内容。（大概率只是写入了缓存，还没有真正写入磁盘）
func (w *BufferedFileWriter) Write(cont []byte) {
	if len(cont) >= len(w.buffer) { //要写的内容比缓存空间还要大，则直接把cont写到文件里去
		w.Flush()
		w.fout.Write(cont)
	} else {
		//先预判buffer能否容下cont
		if w.bufferEndIndex+len(cont) > len(w.buffer) { //不能容下
			w.Flush()
		}
		// append2(w.buffer, w.bufferEndIndex, cont)
		copy(w.buffer[w.bufferEndIndex:], cont) //golang内置的copy函数功能上等价于自己写的append2函数，但比append2函数更高效
		w.bufferEndIndex += len(cont)
	}
}

// 把buffer里的内容全部写入磁盘文件
func (w *BufferedFileWriter) Flush() {
	w.fout.Write(w.buffer[0:w.bufferEndIndex]) //把buffer里的内容写入文件
	w.bufferEndIndex = 0                       //清空buffer
}

// 把src拷贝到dest[index:]里去
func append2(dest []byte, index int, src []byte) {
	for i := 0; i < len(src); i++ {
		dest[index+i] = src[i]
	}
}

// 向文件中写入内容。（大概率只是写入了缓冲，还没有真正写入磁盘）
func (writer *BufferedFileWriter) WriteString(content string) {
	writer.Write([]byte(content))
}

// 直接写文件
func WriteDirect(outFile string) {
	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	for i := 0; i < 100000; i++ {
		fout.WriteString(text)
	}
}

// 带缓冲写文件
func WriteWithBuffer(outFile string) {
	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	writer := NewWriter(fout, 4096)
	defer writer.Flush() //最后，务必把缓冲里残留的内容写入磁盘
	for i := 0; i < 100000; i++ {
		writer.WriteString(text)
	}
}

func main() {
	t1 := time.Now()
	WriteDirect(util.RootPath + "/data/no_buffer.txt")
	t2 := time.Now()
	WriteWithBuffer(util.RootPath + "/data/with_buffer.txt")
	t3 := time.Now()
	fmt.Printf("不用缓冲耗时%dms，用缓冲耗时%dms\n", t2.Sub(t1).Milliseconds(), t3.Sub(t2).Milliseconds())
}

// go run .\io\buffer\

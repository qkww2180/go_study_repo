package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func limitReader() {
	reader := strings.NewReader("daqiaoqiao")
	limitReader := io.LimitReader(reader, 6) //limitReader截取了reader的前6个字节
	content := make([]byte, 100)
	if n, err := limitReader.Read(content); err == nil {
		fmt.Printf("read %s\n", string(content[:n])) //daqiao
	}
	if _, err := limitReader.Read(content); err == io.EOF { //从limitReader里已读不出任何内容
		fmt.Println("no more data available")
	}
}

func multiReader() {
	r1 := strings.NewReader("黄梅时节家家雨\n")
	r2 := strings.NewReader("青草池塘处处蛙\n")
	r3 := strings.NewReader("有约不来过夜半\n")
	r4 := strings.NewReader("闲敲棋子落灯花\n")

	r := io.MultiReader(r1, r2, r3, r4) //注意，这里是有序的
	io.Copy(os.Stdout, r)               //把r流拷贝到标准输出流
	//借助于MultiReader，可以把多个文件合并成一个文件
}

func multiWriter() {
	var (
		writer1 bytes.Buffer
		writer2 bytes.Buffer
	)
	multiWriter := io.MultiWriter(&writer1, &writer2)
	multiWriter.Write([]byte("黄梅时节家家雨\n"))

	fmt.Print(writer1.String()) //黄梅时节家家雨
	fmt.Print(writer2.String()) //黄梅时节家家雨
	//借助于MultiWriter可以把一条日志输出到多个文件里面去
}

func teeReader() {
	var writer bytes.Buffer
	reader := strings.NewReader("黄梅时节家家雨\n")
	teeReader := io.TeeReader(reader, &writer) //从reader里读取的内容既会进入teeReader，也会进入writer
	// io.Copy(os.Stdout, reader)  如果打开此行，则reader里的内容已经被读完了，下面2行不会输出任何内容
	io.Copy(os.Stdout, teeReader) //如果把此行注释掉，则没有任何数据经过teeReader，writer里也不会有任何内容
	fmt.Print(writer.String())
}

func pipeIO() {
	reader, writer := io.Pipe() //writer的内容会直接进入reader，中间没有buffer
	go func() {
		writer.Write([]byte("hello")) //因为中间没有buffer，所以Write操作会阻塞，直到另一个协程准备从reader里读取内容
		writer.Close()
	}()
	content := make([]byte, 100)
	if n, err := reader.Read(content); err == nil {
		fmt.Printf("read %s\n", string(content[:n])) //hello
	}
	reader.Close()
}

func main() {
	limitReader()
	fmt.Println()
	multiReader()
	fmt.Println()
	multiWriter()
	fmt.Println()
	teeReader()
	fmt.Println()
	pipeIO()
}

// go run ./io/advanced_io

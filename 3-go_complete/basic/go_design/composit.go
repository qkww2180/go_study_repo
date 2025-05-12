package main

import "fmt"

type Reader interface {
	Read(int) []byte
}

type Writer interface {
	Write(int) []byte
}

type FileReader struct{}

func (FileReader) Read(n int) []byte {
	fmt.Println("this is FileReader")
	return make([]byte, n)
}

type LogReader struct {
	parent1 Reader //组合
	parent2 Writer
}

func NewLogReader(r Reader) LogReader {
	return LogReader{parent1: r}
}

func (lr LogReader) Read(n int) []byte {
	fmt.Println("this is LogReader")
	array := lr.parent1.Read(n)
	for i := 0; i < n; i++ {
		array[i] /= 2
	}
	return array
}

func read(r Reader) {
	r.Read(8)
}

func readFile(r Reader) {
	r.Read(8)
}

func main2() {
	var fr FileReader = FileReader{}
	read(fr)

	var lr LogReader = LogReader{}
	readFile(lr)
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func readF(fileName string) {
	fin, err := os.Open(fileName) //打开文件，可能发生error
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fin.Close()
	br := bufio.NewReader(fin)

	for {
		if line, err := br.ReadString('\n'); err == nil { //读文件，可能发生error
			fmt.Println(line)
		} else {
			if err == io.EOF { //不同的error，有不同的处理方式
				break
			} else {
				fmt.Println(err)
			}
		}
	}
}

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func appendFile(fout *os.File, infile string) {
	fin, err := os.Open(infile)
	if err != nil {
		log.Panic(err)
	}
	defer fin.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := fin.Read(buffer)
		if err != nil {
			if err == io.EOF {
				if n > 0 {
					fout.Write(buffer[:n])
				}
			} else {
				log.Println(err)
			}
			break
		} else {
			fout.Write(buffer[:n])
		}
	}
}

func main2() {
	mergedFile := "D://Download//合并.mp3"
	fout, err := os.OpenFile(mergedFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	defer fout.Close()

	path := "D://Download//光辉岁月" //把这个目录下的所有文件合并到mergedFile里去
	if fileInfos, err := os.ReadDir(path); err != nil {
		log.Panic(err)
	} else {
		for _, fileInfo := range fileInfos {
			if fileInfo.Type().IsRegular() {
				infile := filepath.Join(path, fileInfo.Name())
				appendFile(fout, infile)
			}
		}
	}
}
